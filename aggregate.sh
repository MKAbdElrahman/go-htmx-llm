#!/bin/bash

# Check if at least two arguments are provided (output file and source directories)
if [ "$#" -lt 2 ]; then
  echo "Usage: $0 <output-file> <source-directory1> [<source-directory2> ...]"
  exit 1
fi

# Output file is the first argument
output_file=$1
shift # Remove the first argument (output file) from the list

# Create or clear the output file
> "$output_file"

# Loop through each source directory provided as argument
for dir in "$@"; do
  # Check if the source directory exists
  if [ -d "$dir" ]; then
    # Loop through each file in the source directory
    for file in "$dir"/*; do
      # If it's a regular file, append its contents to the output file
      if [ -f "$file" ]; then
        cat "$file" >> "$output_file"
        echo -e "\n" >> "$output_file"  # Add a newline after each file's content
        echo "Aggregated content from $file"    
      fi
    done
  else
    echo "Directory $dir does not exist, skipping..."
  fi
done

echo "Aggregation complete. All contents are in $output_file."
