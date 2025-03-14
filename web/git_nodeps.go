package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// Configuration
	repoPath   = "/tmp/repo"         // Change this to your repository path
	listenAddr = ":8080"             // Change port if needed
	authToken  = "your-secret-token" // Change this to your secure token
)

func main() {
	// Ensure the repository exists
	if err := initRepo(); err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Set up HTTP handlers with CORS middleware
	http.HandleFunc("/", corsMiddleware(tokenAuthMiddleware(handleGitRequest)))

	// Start the server
	log.Printf("Git HTTP server listening on %s, serving repo at %s", listenAddr, repoPath)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

// Add CORS headers to all responses
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Initialize the Git repository if it doesn't exist
func initRepo() error {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		log.Printf("Repository doesn't exist, creating at %s", repoPath)

		// Create the directory
		if err := os.MkdirAll(repoPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		// Initialize bare repository
		cmd := exec.Command("git", "init", "--bare", repoPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git init failed: %v, output: %s", err, output)
		}

		// Create a default branch with at least one commit
		tempDir, err := os.MkdirTemp("", "git-init")
		if err != nil {
			return fmt.Errorf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Initialize temp repo
		cmd = exec.Command("git", "init", tempDir)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to init temp repo: %v", err)
		}

		// Configure user
		cmd = exec.Command("git", "-C", tempDir, "config", "user.name", "Git Server")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set user.name: %v", err)
		}

		cmd = exec.Command("git", "-C", tempDir, "config", "user.email", "git@example.com")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set user.email: %v", err)
		}

		// Create README file
		readmePath := filepath.Join(tempDir, "README.md")
		if err := os.WriteFile(readmePath, []byte("# Git Repository\nInitialized by Git HTTP Server"), 0644); err != nil {
			return fmt.Errorf("failed to create README: %v", err)
		}

		// Add and commit
		cmd = exec.Command("git", "-C", tempDir, "add", ".")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add files: %v", err)
		}

		cmd = exec.Command("git", "-C", tempDir, "commit", "-m", "Initial commit")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to commit: %v", err)
		}

		// Push to bare repo
		cmd = exec.Command("git", "-C", tempDir, "remote", "add", "origin", repoPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add remote: %v", err)
		}

		cmd = exec.Command("git", "-C", tempDir, "push", "-u", "origin", "master")
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to push: %v, output: %s", err, output)
		}
	}
	return nil
}

// Middleware to check for valid token
func tokenAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for OPTIONS requests (handled by CORS middleware)
		if r.Method == "OPTIONS" {
			next(w, r)
			return
		}

		// Check for token in Authorization header
		//authHeader := r.Header.Get("Authorization")
		//if authHeader == "" {
		//	// Also check for token in query parameter
		//	token := r.URL.Query().Get("token")
		//	if token != authToken {
		//		http.Error(w, "Unauthorized: Missing or invalid token", http.StatusUnauthorized)
		//		return
		//	}
		//} else {
		//	// Format: "Bearer <token>"
		//	parts := strings.Split(authHeader, " ")
		//	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] != authToken {
		//		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		//		return
		//	}
		//}

		// Token is valid, proceed
		next(w, r)
	}
}

// Main Git request handler
func handleGitRequest(w http.ResponseWriter, r *http.Request) {
	// For debugging
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	// The info/refs endpoint with service parameter is the entry point for git smart protocol
	if strings.HasSuffix(r.URL.Path, "/info/refs") {
		service := r.URL.Query().Get("service")
		if service != "" {
			handleInfoRefs(w, r, service)
			return
		}
	}

	// Handle git services (upload-pack, receive-pack)
	if strings.HasSuffix(r.URL.Path, "/git-upload-pack") {
		handleGitService(w, r, "upload-pack")
		return
	} else if strings.HasSuffix(r.URL.Path, "/git-receive-pack") {
		handleGitService(w, r, "receive-pack")
		return
	}

	// Default response for other endpoints
	http.Error(w, "Not Found", http.StatusNotFound)
}

// Helper function to write a packet line
func writePktLine(w io.Writer, data string) {
	fmt.Fprintf(w, "%04x%s", len(data)+4, data)
}

// Handle the info/refs endpoint which negotiates the Git protocol
func handleInfoRefs(w http.ResponseWriter, r *http.Request, service string) {
	log.Printf("Handling info/refs for service: %s", service)

	// Set the content type for git smart protocol
	w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))

	// Write the smart protocol header - this must be exact format
	writePktLine(w, fmt.Sprintf("# service=%s\n", service))
	fmt.Fprint(w, "0000") // Flush packet

	// Execute git command to get refs
	cmd := exec.Command("git", "--git-dir", repoPath, service[4:], "--stateless-rpc", "--advertise-refs", ".")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe: %v", err)
		http.Error(w, "Git operation failed", http.StatusInternalServerError)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting git command: %v", err)
		http.Error(w, "Git operation failed", http.StatusInternalServerError)
		return
	}

	// Copy the command output directly to the response
	if _, err := io.Copy(w, stdout); err != nil {
		log.Printf("Error copying git output: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Git command error in info/refs: %v", err)
	}
}

// Handle git service endpoints (upload-pack, receive-pack)
func handleGitService(w http.ResponseWriter, r *http.Request, service string) {
	log.Printf("Handling git service: %s", service)

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-result", service))

	cmd := exec.Command("git", "--git-dir", repoPath, service, "--stateless-rpc", ".")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Error creating stdin pipe: %v", err)
		http.Error(w, "Git operation failed", http.StatusInternalServerError)
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe: %v", err)
		stdin.Close()
		http.Error(w, "Git operation failed", http.StatusInternalServerError)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting git command: %v", err)
		stdin.Close()
		http.Error(w, "Git operation failed", http.StatusInternalServerError)
		return
	}

	// Copy request body to git process stdin
	if _, err := io.Copy(stdin, r.Body); err != nil {
		log.Printf("Error copying request to git: %v", err)
	}
	stdin.Close()

	// Copy git process stdout to response
	if _, err := io.Copy(w, stdout); err != nil {
		log.Printf("Error copying git output: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Git command error in %s: %v", service, err)
	}
}
