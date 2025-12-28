const TEMPLATE_PATH = './utils/releasing/templates/Admine-Deploy-Pack' | path expand

def main [version: string, output_path: path, clean?: bool] {
    let do_clean = ($clean | default false)

    setup_tamplate $TEMPLATE_PATH $output_path
    release_vpn_handler $output_path $do_clean
    release_server_handler $output_path $do_clean
    release_bot $output_path $do_clean
    setup_minecraft_server_options $output_path
    create_git_tag $version
}

# Release
def setup_tamplate [template_path: path, output_path: path] {
    if ( ($output_path | path exists)) { error make {msg: "Output path already exists", } }
    cp -r $'($template_path)' $output_path
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
    compile_vpn_handler $clean
    let target_dir = $'($output_path)/vpn_handler'
    if (not ($target_dir | path exists)) { mkdir $target_dir }
    cp ./target/release/vpn_handler $'($target_dir)/vpn_handler'
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
    compile_server_handler $clean
    let target_dir = $'($output_path)/server_handler'
    if (not ($target_dir | path exists)) { mkdir $target_dir }
    cp ./bin/server_handler $'($target_dir)/server_handler'
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
    compile_bot $clean
    let target_dir = $'($output_path)/bot'
    if (not ($target_dir | path exists)) { mkdir $target_dir }
    cp ./dist/bot $'($target_dir)/bot'
}

# Minecraft server
def setup_minecraft_server_options [output_path: path] {
    let minecraft_output_path  = $'($output_path)/minecraft_server'
    if (not ($minecraft_output_path | path exists)) { mkdir $minecraft_output_path }

    # copy variants from repo to output pack
    cp -r ./minecraft_server/* $minecraft_output_path

    cd $minecraft_output_path
    for minecraft_folder in (ls | get name) {
        if ($'./($minecraft_folder)/setup.nu' | path exists) {
            cd $minecraft_folder
            nu setup.nu
            rm -rf *-templates*
            rm -rf setup.nu
            cd ..
        }
    }
}

def create_git_tag [tag_name: string] {
    git tag -a $tag_name -m $'Admine release ($tag_name)'
    git push --tags
}