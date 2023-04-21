mod profiles;

use std::ffi::OsString;

use clap::{command, Parser, Subcommand};

#[derive(Debug, Parser)]
#[command(name = "dfm")]
#[command(about = "A dotfile manager for pair programmers and lazy people.", long_about = None)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Debug, Subcommand)]
enum Commands {
    #[command(external_subcommand)]
    External(Vec<OsString>),
}

fn main() {
    let args = Cli::parse();

    match args.command {
        Commands::External(args) => {
            let plugin_name = format!("dfm-{}", args[0].to_str().unwrap_or_default());
            println!("Calling out to {:?} with {:?}", plugin_name, &args[1..]);
        }
    }
}
