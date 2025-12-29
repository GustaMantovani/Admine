def main [network_id?: string] {
    ls *-templates* | get name | each {|f| 
        let new_name = ($f | str replace "-templates" "")
        cp -r $f $new_name
    }

    touch .env

    if ($network_id != null and ('.env' | path exists) ) {
        let content = $'NETWORK_ID=($network_id)'
        echo $content >> .env
    }
}
