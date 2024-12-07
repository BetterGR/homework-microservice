# Directories
PROTO_DIR = homework_protos
OUTPUT_DIR = homework_protos

# Protobuf compiler command
PROTOC = protoc

# Default target: Compile proto files
all:
	$(PROTOC) -I=$(PROTO_DIR) --go_out=paths=source_relative:$(OUTPUT_DIR) --go-grpc_out=paths=source_relative:$(OUTPUT_DIR) $(PROTO_DIR)/*.proto

# Clean target: Remove generated files
clean:
	rm -f $(OUTPUT_DIR)/*.pb.go

.PHONY: all clean
