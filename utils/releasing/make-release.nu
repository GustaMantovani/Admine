const TEMPLATE_PATH = './utils/releasing/templates/Admine-Deploy-Pack' | path expand

def main [
    version: string
    output_path?: path = './releases'
    --clean = false
    --force = false
    --dev = false
    --push_tags = true
] {
    let repo_root = pwd | path expand
    let output_path = $output_path | path expand
    let $output_path  = $'($output_path)/($version)/admine-deploy-pack-linux-x86_64-($version)'

    print $"🚀 Starting Admine release ($version)"

    setup_template $TEMPLATE_PATH $output_path $force

    print "📦 Building VPN Handler..."
    release_vpn_handler $output_path $clean

    print "📦 Building Server Handler..."
    release_server_handler $output_path $clean

    print "📦 Building Bot..."
    release_bot $output_path $clean

    if ($dev) { copy_dev_configs $repo_root $output_path }

    print "🗜️  Creating archives..."
    create_compress_archives $output_path

    if (not $dev) { create_git_tag $version $push_tags $force }

    print $"✅ Release ($version) completed successfully!"
}

# Release
def setup_template [template_path: path, output_path: path, force: bool] {
    if ($output_path | path exists) {
        if (not $force) { 
            error make { msg: "Output path already exists" } 
        } else {
            print "⚠️  Warning: Removing existing output directory..."
            mkdir $output_path
        }
    } else {
        mkdir $output_path
    }

    print "📋 Copying template files..."
    cd $template_path
    ls | par-each { |row| 
        cp -r $'($template_path)/($row.name)' $output_path 
    }
    print "  ✓ Template files copied"
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

# Dev mode: overwrite release configs with the local project configs so the
# build can be tested immediately without manual reconfiguration.
# Must run before create_compress_archives so the files are included in the archives.
def copy_dev_configs [repo_root: path, output_path: path] {
    print "🔧 Dev mode: copying local configs into release..."

    let configs = [
        [$'($repo_root)/server_handler/server_handler_config.yaml', $'($output_path)/server_handler/server_handler_config.yaml'],
        [$'($repo_root)/vpn_handler/etc/vpn_handler_config.toml',   $'($output_path)/vpn_handler/etc/vpn_handler_config.toml'],
        [$'($repo_root)/bot/bot_config.json',                        $'($output_path)/bot/bot_config.json'],
    ]

    for pair in $configs {
        let src = $pair.0
        let dst = $pair.1
        if ($src | path exists) {
            cp $src $dst
            print $"  ✓ ($src | path basename)"
        } else {
            print $"  ⚠️  Skipped \(not found\): ($src | path relative-to $repo_root)"
        }
    }
}

def create_compress_archives [output_path: path] {
    let parent_dir = $output_path | path dirname
    let folder_name = $output_path | path basename
    let tar_file = $'($folder_name).tar.gz'
    let zip_file = $'($folder_name).zip'
    
    cd $parent_dir
    
    print "  📦 Creating tar.gz..."
    tar -czf $tar_file $folder_name

    print "  📦 Creating zip..."
    ^zip -r $zip_file $folder_name

    print "  ✓ Archives created"
}

def create_git_tag [tag_name: string, push_tags: bool, force: bool] {
    if (not $force) {
        if (git tag -l $tag_name | is-not-empty) {
            error make { msg: $"Tag ($tag_name) already exists" }
        }

        git fetch --tags
        if (git ls-remote --tags origin $tag_name | is-not-empty) {
            error make { msg: $"Tag ($tag_name) already exists on remote" }
        }
    } else {
        print "⚠️  Force enabled: checking existing tags..."
        git fetch --tags

        if (git tag -l $tag_name | is-not-empty) {
            print $"  ⚠️  Deleting existing local tag ($tag_name)..."
            git tag -d $tag_name
        }

        if (git ls-remote --tags origin $tag_name | is-not-empty) {
            print $"  ⚠️  Deleting existing remote tag ($tag_name)..."
            git push origin $':refs/tags/($tag_name)'
        }
    }

    let current_branch = (git branch --show-current | str trim)
    if $current_branch not-in ["main" "master"] {
        error make { 
            msg: $"You can only create tags from main/master branch. Current branch: ($current_branch)" 
        }
    }

    git tag -a $tag_name -m $'Admine ($tag_name)'

    if ($push_tags) {
        git push origin $tag_name
    }
}
