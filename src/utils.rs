use std::path::PathBuf;
use std::str::FromStr;

pub fn expand_path(path: &str) -> PathBuf {
    let expanded = shellexpand::tilde(path);
    PathBuf::from_str(&expanded).unwrap_or_else(|_| panic!("invalid path: {}", path))
}
