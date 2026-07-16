use std::env;

fn grep() {
    let args: Vec<_> = env::args().collect();

    if args.len() < 3 {
        println!("Not enough arguments");
        return;
    }

    let text = &args[1];
    let file_path = &args[2];
}
