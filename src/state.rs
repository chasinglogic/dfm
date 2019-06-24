use crate::util;
use serde::{Deserialize, Serialize};
use std::error::Error;
use std::fs::File;
use std::io;
use std::io::{Read, Write};
use std::process;

#[derive(Deserialize, Serialize)]
pub struct State {
    pub current_profile: String,
}

impl State {
    pub fn default() -> State {
        State {
            current_profile: "".to_string(),
        }
    }

    pub fn load() -> Option<State> {
        let sf = util::state_file_p();
        let mut contents = String::new();
        match File::open(sf) {
            Ok(mut f) => {
                if let Err(_) = f.read_to_string(&mut contents) {
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

    pub fn save(self) -> Result<(), io::Error> {
        let sf = util::state_file_p();
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
