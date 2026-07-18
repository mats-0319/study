mod grep;
mod guess_number;
mod practice_8_1;
mod practice_8_3;

fn practice() {
    // practice_8_1::f1(&vec![3, 5, 2, 1, 4, 5, 3]);
    // practice_8_3::f3();
}

fn demo() {
    // guess_number::guess_number();
    grep::grep();
}

struct MyBox<T>(T);

impl<T> MyBox<T> {
    fn new(x: T) -> MyBox<T> {
        MyBox(x)
    }
}

use std::ops::Deref;

impl<T> Deref for MyBox<T> {
    type Target = T;

    fn deref(&self) -> &Self::Target {
        &self.0
    }
}

fn hello(name: &str) {
    println!("Hello, {name}!");
}

fn main() {
    // practice();
    // demo();

    let a = 5;
    let b = MyBox::new(String::from("Rust"));

    hello(&b);
}
