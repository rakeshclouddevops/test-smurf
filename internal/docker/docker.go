package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/fatih/color"
	"github.com/moby/term"
	"github.com/pterm/pterm"
)

// BuildOptions struct to hold options for Docker build
type BuildOptions struct {
	DockerfilePath string
	NoCache        bool
	BuildArgs      map[string]*string
	Target         string
}

// createTarArchive creates a tar archive of the entire build context directory.
func createTarArchive(contextDir string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	err := filepath.Walk(contextDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(file)
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.Mode().IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			defer data.Close()
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

// Build builds a Docker image from a specified Dockerfile.
func Build(imageName, tag string, opts BuildOptions) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	contextDir := filepath.Dir(opts.DockerfilePath)
	tarStream, err := createTarArchive(contextDir)
	if err != nil {
		return err
	}

	options := types.ImageBuildOptions{
		Tags:        []string{fmt.Sprintf("%s:%s", imageName, tag)},
		Dockerfile:  filepath.Base(opts.DockerfilePath),
		NoCache:     opts.NoCache,
		BuildArgs:   opts.BuildArgs,
		Target:      opts.Target,
		Remove:      true,
		ForceRemove: true,
		PullParent:  true,
	}

	spinner, _ := pterm.DefaultSpinner.Start("Building Docker image...")
	buildResponse, err := cli.ImageBuild(ctx, tarStream, options)
	if err != nil {
		spinner.Fail("Failed to start the build process")
		color.New(color.FgRed).Println(err)
		return err
	}

	defer buildResponse.Body.Close()

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	err = jsonmessage.DisplayJSONMessagesStream(buildResponse.Body, os.Stdout, termFd, isTerm, nil)
	if err != nil {
		spinner.Fail("Failed during the build process")
		color.New(color.FgRed).Println(err)
		return err
	}

	spinner.Success("Docker image built successfully")
	color.New(color.FgGreen).Printf("Successfully built %s:%s\n", imageName, tag)

	return nil
}

// TagOptions struct to hold options for tagging a Docker image
type TagOptions struct {
	Source string
	Target string 
}

// TagImage tags a local Docker image for use in a remote repository.
func TagImage(opts TagOptions) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		color.New(color.FgRed).Printf("Error creating Docker client: %v\n", err)
		return err
	}

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Tagging image %s as %s...", opts.Source, opts.Target))
	if err := cli.ImageTag(ctx, opts.Source, opts.Target); err != nil {
		spinner.Fail(fmt.Sprintf("Failed to tag image: %v", err))
		color.New(color.FgRed).Println(err)
		return err
	}

	spinner.Success(fmt.Sprintf("Successfully tagged %s as %s", opts.Source, opts.Target))
	color.New(color.FgGreen).Printf("Image tagged successfully: %s -> %s\n", opts.Source, opts.Target)
	return nil
}

// PushOptions struct to hold options for pushing a Docker image
type PushOptions struct {
	ImageName string 
}

// PushImage pushes a Docker image to a Docker registry.
func PushImage(opts PushOptions) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		pterm.Error.Println("Error creating Docker client:", err)
		return err
	}

	authConfig := registry.AuthConfig{
		Username: os.Getenv("DOCKER_USERNAME"),
		Password: os.Getenv("DOCKER_PASSWORD"),
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		pterm.Error.Println("Error encoding auth config:", err)
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Pushing image %s...", opts.ImageName))
	options := image.PushOptions{
		RegistryAuth: authStr,
	}

	responseBody, err := cli.ImagePush(ctx, opts.ImageName, options)
	if err != nil {
		spinner.Fail("Failed to push the image: " + err.Error())
		return err
	}
	defer responseBody.Close()

	return handleDockerResponse(responseBody, spinner, opts)
}

func handleDockerResponse(responseBody io.ReadCloser, spinner *pterm.SpinnerPrinter, opts PushOptions) error {
	decoder := json.NewDecoder(responseBody)
	var lastProgress int
	for {
		var msg jsonmessage.JSONMessage
		if err := decoder.Decode(&msg); err == io.EOF {
			break
		} else if err != nil {
			pterm.Error.Println("Error decoding JSON:", err)
			return err
		}

		if msg.Error != nil {
			pterm.Error.Println("Error from Docker:", msg.Error.Message)
			return fmt.Errorf(msg.Error.Message)
		}

		if msg.Progress != nil && msg.Progress.Total > 0 {
			current := int(msg.Progress.Current * 100 / msg.Progress.Total)
			if current > lastProgress {
				progressMessage := fmt.Sprintf("Pushing image %s... %d%%", opts.ImageName, current)
				spinner.UpdateText(progressMessage)
				fmt.Printf("\r%s", pterm.Green(progressMessage)) 
				lastProgress = current
			}
		}

		if msg.Stream != "" {
			fmt.Print(pterm.Blue(msg.Stream)) 
		}
	}

	spinner.Success("Image push complete.")
	pterm.Success.Println("Successfully pushed image:", opts.ImageName)
	return nil
}

// Scout scans a Docker image for known vulnerabilities using 'docker scout cves'.
func Scout(dockerTag, sarifFile string) error {
	ctx := context.Background()

	args := []string{"scout", "cves", dockerTag}

	if sarifFile != "" {
		args = append(args, "--output", sarifFile)
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	spinner, _ := pterm.DefaultSpinner.Start("Running 'docker scout cves'")
	defer spinner.Stop()

	err := cmd.Run()

	spinner.Stop()

	outStr := stdoutBuf.String()
	errStr := stderrBuf.String()

	if err != nil {
		pterm.Error.Println("Error running 'docker scout cves':", err)
		if errStr != "" {
			pterm.Error.Println(errStr)
		}
		return fmt.Errorf("failed to run 'docker scout cves': %w", err)
	}

	if outStr != "" {
		pterm.Info.Println("Docker Scout CVEs output:")
		fmt.Println(color.YellowString(outStr))
	}

	if sarifFile != "" {
		if _, err := os.Stat(sarifFile); err == nil {
			pterm.Success.Println("SARIF report saved to:", sarifFile)
		} else {
			pterm.Warning.Println("Expected SARIF report not found at:", sarifFile)
		}
	}

	pterm.Success.Println("Scan completed successfully.")
	return nil
}
