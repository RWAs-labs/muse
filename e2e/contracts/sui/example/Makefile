.PHONY: clean build all

# Default target
all: clean build

# Clean build directory
clean:
	rm -rf build/

# Build the package and generate bytecode
build:
	sui move build
	cp build/example/bytecode_modules/token.mv ../token.mv	
	cp build/example/bytecode_modules/connected.mv ../connected.mv

# Help target
help:
	@echo "Available targets:"
	@echo "  all    - Clean and build everything (default)"
	@echo "  clean  - Remove build directory"
	@echo "  build  - Build the package and generate bytecode"
	@echo "  help   - Show this help message" 