use std::io;
#[cfg(target_family = "unix")]
use std::os::unix::fs::symlink;
#[cfg(target_os = "windows")]
use std::os::windows::fs::symlink_file as symlink;
use std::path::{Path, PathBuf};

#[derive(Debug, PartialEq, Eq, Clone)]
pub struct Info {
    src: PathBuf,
    dst: PathBuf,
    profile_dir: PathBuf,
}

impl Info {
    pub fn new(src: &Path, profile_dir: &Path, target_dir: &Path) -> Info {
        let new_src = src.to_path_buf();
        let mut i = Info {
            src: new_src.clone(),
            profile_dir: profile_dir.to_path_buf(),
            // will be overwritten by retarget below
            dst: new_src,
        };
        i.retarget(target_dir);
        i
    }

    pub fn retarget(&mut self, target_dir: &Path) {
        let filename = self.src.strip_prefix(&self.profile_dir).unwrap();
        let mut dst = target_dir.to_path_buf();
        dst.push(filename);
        self.dst = dst;
    }

    pub fn link(&self) -> Result<(), io::Error> {
        symlink(&self.src, &self.dst)
    }

    pub fn get_dst(&self) -> PathBuf {
        self.dst.clone()
    }
}

#[cfg(test)]
pub mod test {
    use super::Info;
    use std::path::{Path, PathBuf};

    #[test]
    fn test_new() {
        let src = Path::new("/home/foo/.config/dfm/profile/bar/.bashrc");
        let profile_dir = Path::new("/home/foo/.config/dfm/profile/bar");
        let target_dir = Path::new("/home/foo");
        let info = Info::new(src, profile_dir, target_dir);
        assert_eq!(
            info,
            Info {
                src: src.to_path_buf(),
                dst: PathBuf::from("/home/foo/.bashrc"),
                profile_dir: profile_dir.to_path_buf(),
            }
        )
    }
}
