use std::env;
use std::fs::create_dir_all;
use std::path::{Path, PathBuf};
use std::process;

use crate::profile::Profile;

fn ensure_exists(p: &Path) {
    if !p.exists() {
        create_dir_all(p).expect(&format!("Unable to create directory: {}", p.display()))
    }
}

pub fn xdg_dir() -> PathBuf {
    match env::var("XDG_CONFIG_HOME") {
        Ok(path) => PathBuf::from(path),
        Err(_) => {
            let home = env::var("HOME").unwrap_or("".to_string());
            let mut home_p = PathBuf::from(home);
            home_p.push(".config");
            home_p
        }
    }
}

pub fn dfm_dir() -> PathBuf {
    let p = match env::var("DFM_CONFIG_DIR") {
        Ok(path) => PathBuf::from(path),
        Err(_) => {
            let mut path = xdg_dir();
            path.push("dfm");
            path
        }
    };
    ensure_exists(&p);
    p
}

pub fn state_file_p() -> PathBuf {
    let mut p = dfm_dir();
    p.push("state.yml");
    p
}

pub fn profile_storage_dir() -> PathBuf {
    let mut p = dfm_dir();
    p.push("profiles");
    ensure_exists(&p);
    p
}

pub fn profile_dir(name: &str) -> PathBuf {
    let mut storage = profile_storage_dir();
    storage.push(name);
    storage
}

pub fn load_profile(name: &str) -> Profile {
    if name == "" {
        println!("No current profile and no profile provided.");
        println!("Try running dfm link <profile name> to set an active profile.");
        process::exit(1);
    }

    let dir = profile_dir(name);
    if !dir.exists() {
        println!("Profile directory does not exist {}", dir.display());
        println!("Cannot load from non-existent directory");
        process::exit(1);
    }

    // TODO: handle error better
    Profile::load(&dir).expect(&format!("Unable to load profile: {}", name))
}
