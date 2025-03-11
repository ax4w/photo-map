#!/bin/bash
# This script iterates over the "images" directory and creates thumbnails in the "thumbs" directory.
# It preserves the directory structure and uses ImageMagick's convert command to resize images.

SRC_DIR="images"
DST_DIR="thumbs"
THUMB_SIZE="300x300"  # Adjust size as needed

# Create the destination directory if it doesn't exist
mkdir -p "$DST_DIR"

# Iterate over supported image files in the SRC_DIR
find "$SRC_DIR" -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" -o -iname "*.gif" -o -iname "*.webp" \) | while IFS= read -r file; do
    # Compute the relative path of the image file
    rel_path="${file#$SRC_DIR/}"
    # Determine the destination file path in the DST_DIR with the same relative path
    dst_file="$DST_DIR/$rel_path"
    # Create the destination directory if it doesn't exist
    mkdir -p "$(dirname "$dst_file")"
    # Resize the image and save it to the destination
    convert "$file" -resize "$THUMB_SIZE" "$dst_file"
    echo "Thumbnail created: $dst_file"
done
