use std::env;
use std::fs::create_dir_all;
use std::path::{Path, PathBuf};
use std::process;

use crate::profile::Profile;

pub fn ensure_exists(p: &Path) {
    if !p.exists() {
        if let Err(e) = create_dir_all(p) {
            println!("Unable to create directory {}: {}", p.display(), e);
            process::exit(1);
        }
    }
}

pub fn xdg_dir() -> PathBuf {
    match env::var("XDG_CONFIG_HOME") {
        Ok(path) => PathBuf::from(path),
        Err(_) => {
            let home = env::var("HOME").unwrap_or_default();
            let mut home_p = PathBuf::from(home);
            home_p.push(".config");
            home_p
        }
    }
}

pub fn cfg_dir(cfd: Option<&Path>) -> PathBuf {
    match cfd {
        Some(dir) => dir.to_path_buf(),
        None => match env::var("DFM_CONFIG_DIR") {
            Ok(path) => PathBuf::from(path),
            Err(_) => xdg_dir(),
        },
    }
}

pub fn dfm_dir(cfd: &Path) -> PathBuf {
    let mut p = cfg_dir(Some(cfd));
    p.push("dfm");
    ensure_exists(&p);
    p
}

pub fn state_file_p(cfd: &Path) -> PathBuf {
    let mut p = dfm_dir(cfd);
    p.push("state.yml");
    p
}

pub fn profile_storage_dir(cfd: &Path) -> PathBuf {
    let mut p = dfm_dir(cfd);
    p.push("profiles");
    ensure_exists(&p);
    p
}

pub fn profile_dir(name: &str, cfd: &Path) -> PathBuf {
    let mut storage = profile_storage_dir(cfd);
    storage.push(name);
    storage
}

pub fn load_profile(name: &str, cfd: &Path) -> Profile {
    if name == "" {
        println!("No current profile and no profile provided.");
        println!("Try running dfm link <profile name> to set an active profile.");
        process::exit(1);
    }

    let dir = profile_dir(name, cfd);
    if !dir.exists() {
        println!("Profile directory does not exist {}", dir.display());
        println!("Cannot load from non-existent directory");
        process::exit(1);
    }

    match Profile::load(&dir) {
        Ok(p) => p,
        Err(e) => {
            println!("Unable to load profile {}: {}", name, e);
            process::exit(1);
        }
    }
}
