#[derive(Debug, serde::Deserialize, serde::Serialize)]
pub struct State {
    current_profile: String,
}

impl Default for State {
    fn default() -> Self {
        State {
            current_profile: "".to_string(),
        }
    }
}

impl State {
    pub fn load(fp: &Path) -> Result<State, io::Error> {
        let fh = File::open(fp)?;
        let buffer = BufReader::new(fh);
        Ok(serde_json::from_reader(buffer)?)
    }

    pub fn save(&self, filepath: &Path) -> Result<(), io::Error> {
        if let Some(parent) = filepath.parent() {
            if !parent.exists() {
                fs::create_dir_all(parent).expect("Unable to create dfm directory!");
            }
        }

        let file_handle = File::create(filepath)?;
        Ok(serde_json::to_writer(file_handle, self)?)
    }
}

pub fn home_dir() -> PathBuf {
    let home = env::var("HOME").unwrap_or("".to_string());
    PathBuf::from(home)
}

pub fn dfm_dir() -> PathBuf {
    let mut path = home_dir();
    path.push(".config");
    path.push("dfm");
    path
}

pub fn state_file() -> PathBuf {
    let mut state_fp = dfm_dir();
    state_fp.push("state.json");
    state_fp
}

pub fn profiles_dir() -> PathBuf {
    let mut path = dfm_dir();
    path.push("profiles");
    if !path.exists() {
        fs::create_dir_all(&path).expect("Unable to create profiles directory!");
    }

    path
}

pub fn load_profile(name: &str) -> Profile {
    let mut path = profiles_dir();
    path.push(name);
    Profile::load(&path)
}

pub fn force_available(profile: Option<Profile>) -> Profile {
    match profile {
        None => {
            eprintln!("No profile is currently loaded!");
            process::exit(1);
        }
        Some(p) => p,
    }
}

pub fn load_or_default() -> Result<State, io::Error> {
    let state_fp = state_file();
    match State::load(&state_fp) {
        Ok(state) => state,
        Err(err) => match err.kind() {
            io::ErrorKind::NotFound => State::default(),
            _ => panic!("{}", err),
        },
    };
}
