const TEMPLATE_PATH = './utils/releasing/templates/Admine-Deploy-Pack' | path expand

def main [version: string, output_path: path, clean?: bool] {
    let do_clean = ($clean | default false)
    let output_path = $output_path | path expand
    let $output_path  = $'($output_path)/admine-deploy-pack-linux-x86_64-($version)'

    print $"🚀 Starting Admine release ($version)"

    setup_tamplate $TEMPLATE_PATH $output_path
    
    print "📦 Building VPN Handler..."
    release_vpn_handler $output_path $do_clean
    
    print "📦 Building Server Handler..."
    release_server_handler $output_path $do_clean
    
    print "📦 Building Bot..."
    release_bot $output_path $do_clean
    
    print "🗜️  Creating archives..."
    create_compress_archives $output_path $do_clean
    
    create_git_tag $version
    print $"✅ Release ($version) completed successfully!"
}

# Release
def setup_tamplate [template_path: path, output_path: path] {
    if ($output_path | path exists) {
        error make { msg: "Output path already exists" }
    } else {
        mkdir $output_path
    }

    print "📋 Copying template files..."
    cd $template_path
    ls | par-each { |row| 
        cp -r $'($template_path)/($row.name)' $output_path 
    }
}

# VPN Handler
def compile_vpn_handler [clean: bool] {
    if ($clean) {
        cargo clean
    }
    cargo build --release
}

def release_vpn_handler [output_path: path, clean: bool] {
    cd ./vpn_handler
    
    print "  ⚙️  Compiling Rust binary..."
    compile_vpn_handler $clean
    
    let target_dir = $'($output_path)/vpn_handler'
    if (not ($target_dir | path exists)) { mkdir $target_dir }
    cp ./target/release/vpn_handler $'($target_dir)/vpn_handler'
    print "  ✓ VPN Handler ready"
}

# Server Handler
def compile_server_handler [clean: bool] {
    if ($clean) {
        make clean
    }
    make build
}

def release_server_handler [output_path: path, clean: bool] {
    cd ./server_handler
    let binary_path = './bin/server_handler'
    let target_dir = $'($output_path)/server_handler'

    if (not ($binary_path | path exists) or $clean) {
        print "  ⚙️  Compiling Go binary..."
        compile_server_handler $clean    
    }

    if (not ($target_dir | path exists)) { mkdir $target_dir }
    cp $binary_path $'($target_dir)/server_handler'
    print "  ✓ Server Handler ready"
}

# Bot
def compile_bot [clean: bool] {
    if ($clean) {
        make clean
    }
    make build-release
}

def release_bot [output_path: path, clean: bool] {
    cd ./bot
    let binary_path = './dist/bot'
    let target_dir = $'($output_path)/bot'

    if (not ($binary_path | path exists) or $clean) {
        print "  ⚙️  Building Python binary..."
        compile_bot $clean
    }

    if (not ($target_dir | path exists)) { mkdir $target_dir }
    cp $binary_path $'($target_dir)/bot'
    print "  ✓ Bot ready"
}

def create_compress_archives [output_path: path, clean: bool] {
    let parent_dir = $output_path | path dirname
    let folder_name = $output_path | path basename
    let tar_file = $'($folder_name).tar.gz'
    let zip_file = $'($folder_name).zip'
    
    cd $parent_dir
    
    if (not ($tar_file | path exists) or $clean) { 
        print "  📦 Creating tar.gz..."
        tar -czf $tar_file $folder_name 
    }
    if (not ($zip_file | path exists) or $clean) { 
        print "  📦 Creating zip..."
        ^zip -r $zip_file $folder_name 
    }
    print "  ✓ Archives created"
}

def create_git_tag [tag_name: string] {
    git tag -a $tag_name -m $'Admine ($tag_name)'
}