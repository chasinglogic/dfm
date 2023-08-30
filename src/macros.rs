use lazy_static::lazy_static;
use std::env;

lazy_static! {
    pub static ref VERBOSITY_ENABLED: bool = env::var("DFM_VERBOSE").is_ok();
}

#[macro_export]
macro_rules! debug {
    ($($arg:tt)*) => {
        if *$crate::macros::VERBOSITY_ENABLED {
            println!($($arg)*);
        }
    };
}
