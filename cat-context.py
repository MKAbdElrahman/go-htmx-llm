import os

# Map file extensions to Markdown code block languages
EXTENSION_TO_LANGUAGE = {
    ".go": "go",
    ".py": "python",
    ".js": "javascript",
    ".java": "java",
    ".html": "html",
    ".css": "css",
    ".md": "markdown",
    ".txt": "text",
    # Add more extensions and languages as needed
}

def concatenate_files(inputs, output_file, ignore_extensions=None):
    if ignore_extensions is None:
        ignore_extensions = []  # Default to an empty list if no extensions are provided

    with open(output_file, 'w') as outfile:
        for item in inputs:
            if os.path.isfile(item):
                # If the item is a file, check its extension
                if any(item.endswith(ext) for ext in ignore_extensions):
                    print(f"Ignoring file (extension excluded): {item}")
                    continue
                try:
                    with open(item, 'r') as infile:
                        # Get file extension and corresponding language
                        file_extension = os.path.splitext(item)[1]
                        language = EXTENSION_TO_LANGUAGE.get(file_extension, "")
                        
                        # Write metadata and content in Markdown format
                        outfile.write(f"## File: `{item}`\n")
                        outfile.write(f"- **Type**: File\n")
                        outfile.write(f"- **Extension**: `{file_extension}`\n")
                        outfile.write(f"### Content:\n")
                        
                        # Add language-specific code block syntax
                        if language:
                            outfile.write(f"```{language}\n")
                        else:
                            outfile.write("```\n")
                        
                        outfile.write(infile.read())
                        outfile.write("\n```\n\n")  # Add extra spacing between files
                except Exception as e:
                    print(f"Error reading {item}: {e}")
            elif os.path.isdir(item):
                # If the item is a directory, recursively process all files in it
                for root, _, files in os.walk(item):
                    for file in files:
                        file_path = os.path.join(root, file)
                        if any(file.endswith(ext) for ext in ignore_extensions):
                            print(f"Ignoring file (extension excluded): {file_path}")
                            continue
                        try:
                            with open(file_path, 'r') as infile:
                                # Get file extension and corresponding language
                                file_extension = os.path.splitext(file)[1]
                                language = EXTENSION_TO_LANGUAGE.get(file_extension, "")
                                
                                # Write metadata and content in Markdown format
                                outfile.write(f"## File: `{file_path}`\n")
                                outfile.write(f"- **Type**: File\n")
                                outfile.write(f"- **Extension**: `{file_extension}`\n")
                                outfile.write(f"### Content:\n")
                                
                                # Add language-specific code block syntax
                                if language:
                                    outfile.write(f"```{language}\n")
                                else:
                                    outfile.write("```\n")
                                
                                outfile.write(infile.read())
                                outfile.write("\n```\n\n")  # Add extra spacing between files
                        except Exception as e:
                            print(f"Error reading {file_path}: {e}")
            else:
                print(f"Invalid input: {item} is neither a file nor a directory.")

if __name__ == "__main__":
    # List of file paths or directories to concatenate
    inputs = [
        "./chat",
        "./promptprocessing",
        "./cmd", 
        "./index.html",
        "./prompt.md",
    ]

    # Output file where the concatenated content will be written
    output_file = "context.md"  # Save as a Markdown file

    # List of file extensions to ignore (e.g., [".log", ".tmp"])
    ignore_extensions = [".log", ".tmp"]

    concatenate_files(inputs, output_file, ignore_extensions)
    print(f"All files have been concatenated into {output_file}")