use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use std::thread;
use std::time::Duration;

fn main() {
    // Create an atomic boolean to track if the application should stop running
    let running = Arc::new(AtomicBool::new(true));
    let r = running.clone();

    // Handle SIGTERM (and other termination signals)
    ctrlc::set_handler(move || {
        r.store(false, Ordering::SeqCst);
    }).expect("Error setting Ctrl-C handler");

    // Main loop
    while running.load(Ordering::SeqCst) {
        println!("I'm running!");
        thread::sleep(Duration::from_secs(1));
    }
}