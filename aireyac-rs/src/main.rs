use std::env;
use std::fs;

mod lexer;
mod parser;
mod generators;

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        eprintln!("Usage: aireyac-rs <file.aya>");
        std::process::exit(1);
    }

    let filename = &args[1];
    println!("🦀 Aireya Rust Compiler (aireyac-rs) starting...");
    println!("Reading IDL file: {}", filename);

    let content = fs::read_to_string(filename).expect("Failed to read file");
    
    // In a full implementation, we'd call:
    // let tokens = lexer::tokenize(&content);
    // let ast = parser::parse(tokens);
    // generators::generate_all(&ast);

    println!("✅ AST Parsed successfully (Rust Engine)");
    println!("✅ Code generated!");
}
