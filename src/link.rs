use std::io;
use std::fs::{remove_file, symlink_metadata};
use std::os::unix::fs::symlink;
use std::path::{Path, PathBuf};

#[derive(Debug, PartialEq, Eq, Clone)]
pub struct Info {
    pub src: PathBuf,
    pub dst: PathBuf,
    pub overwrite: bool,
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
            overwrite: false,
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
        match symlink_metadata(&self.dst) {
            Ok(metadata) => {
                if metadata.file_type().is_file() && !self.overwrite {
                    println!("File {} already exists, refusing to remove without --overwrite", self.dst.display());
                    // While we want to report to the user we did
                    // nothing we don't want to stop execution of the
                    // program so report an Ok result.
                    return Ok(())
                }

                remove_file(&self.dst)?;
            },
            Err(e) => match e.kind() {
                // if the user doesn't have read access to the file
                // this should be fatal
                io::ErrorKind::PermissionDenied => {
                    return Err(e);
                },
                _ => (),
            }
        };

        match symlink(&self.src, &self.dst) {
            Ok(()) => Ok(()),
            Err(e) => match e.kind() {
                _ => Err(e), // TODO: handle errors
            },
        }
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
                overwrite: false,
            }
        )
    }
}
