.PHONY: all clean

# Define platforms
PLATFORMS := windows_amd64 linux_amd64 android_arm64
RELEASE_FOLDER := release

# Main target
all: $(addprefix $(RELEASE_FOLDER)/TempFiles_api_, $(PLATFORMS))

# Rule to build binaries for each platform
$(RELEASE_FOLDER)/TempFiles_api_%:
	@echo "Building binary for $*..."
	CGO_ENABLED=0 GOOS=$(word 1,$(subst _, ,$*)) GOARCH=$(word 2,$(subst _, ,$*)) go build -tags netgo -ldflags '-s -w' -o "$@" . 
	@echo "Binary for $* built successfully."

# Clean target
clean:
	@echo "Cleaning up..."
	@rm -rf $(RELEASE_FOLDER)
	@echo "Cleanup complete."
