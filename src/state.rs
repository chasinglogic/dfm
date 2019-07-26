use std::error::Error;
use std::fs::File;
use std::io;
use std::io::{Read, Write};
use std::path::Path;
use std::process;

use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Debug)]
pub struct State {
    pub current_profile: String,
}

impl State {
    pub fn default() -> State {
        State {
            current_profile: "".to_string(),
        }
    }

    pub fn load(sf: &Path) -> Option<State> {
        let mut contents = String::new();
        match File::open(sf) {
            Ok(mut f) => {
                if f.read_to_string(&mut contents).is_err() {
                    return None;
                };

                match serde_yaml::from_str(&contents) {
                    Ok(s) => Some(s),
                    Err(_) => None,
                }
            }
            Err(_) => None,
        }
    }

    pub fn save(self, sf: &Path) -> Result<(), io::Error> {
        let display = sf.display();

        let mut file = match File::create(&sf) {
            Err(why) => {
                println!(
                    "Couldn't save state file {}: {}",
                    display,
                    why.description()
                );
                process::exit(1);
            }
            Ok(file) => file,
        };

        let bytes = serde_yaml::to_vec(&self).expect("Failed to serialize app state");
        file.write_all(&bytes)
    }
}
