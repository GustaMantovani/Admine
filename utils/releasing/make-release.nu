const TAMPLATE_PATH = '.utils/releasing/templates/Admine-Deploy-Pack'

def main [version: string, output_path: path, clean?: bool] {

}

#Release
def setup_tamplate [template_path, output_path: path] {
    
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
    cp ./target/release/vpn_handler $'($output_path)/vpn_handler'
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
    compile_vpn_handler $clean
    cp ./bin/server_handler $'($output_path)/server_handler'
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
    compile_vpn_handler $clean
    cp ./dist/bot $'($output_path)/bot'
}

# Minecraft server
def setup_minecraft_server_options [] {
    
}