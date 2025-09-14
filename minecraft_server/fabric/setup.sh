for file in *-template*; do
    new_name="${file/-templates/}"
    cp -r "$file" "$new_name"
done
