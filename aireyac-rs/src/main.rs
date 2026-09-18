use std::env;
use std::fs;

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        eprintln!("Usage: aireyac-rs <file.aya>");
        std::process::exit(1);
    }

    let filename = &args[1];
    println!("🦀 Aireya Rust Compiler (aireyac-rs) starting...");

    let content = fs::read_to_string(filename).expect("Failed to read file");

    let tokens = tokenize(&content);
    let file = parse(&tokens);

    println!("✅ AST Parsed: {} structs, {} services", file.structs.len(), file.services.len());
    println!("   Namespace: {}", file.namespace);
    for s in &file.structs {
        println!("   Struct: {} ({} fields)", s.name, s.fields.len());
    }
    for s in &file.services {
        println!("   Service: {} ({} RPCs)", s.name, s.rpcs.len());
    }
    println!("✅ Compilation complete.");
}

// ─── Token ────────────────────────────────────────────────────
#[derive(Debug, Clone, PartialEq)]
enum TokenType {
    Namespace, Struct, Service, Rpc, Enum, Stream,
    LBrace, RBrace, LParen, RParen, Colon, Semicolon, Arrow, Comma, Eq, Dot,
    Ident(String), Number(String), Eof,
}

// ─── Lexer ────────────────────────────────────────────────────
fn tokenize(src: &str) -> Vec<TokenType> {
    let mut tokens = Vec::new();
    let chars: Vec<char> = src.chars().collect();
    let mut i = 0;
    while i < chars.len() {
        let c = chars[i];
        if c.is_whitespace() { i += 1; continue; }
        if c == '/' && i + 1 < chars.len() && chars[i + 1] == '/' {
            while i < chars.len() && chars[i] != '\n' { i += 1; }
            continue;
        }
        match c {
            '{' => { tokens.push(TokenType::LBrace); i += 1; }
            '}' => { tokens.push(TokenType::RBrace); i += 1; }
            '(' => { tokens.push(TokenType::LParen); i += 1; }
            ')' => { tokens.push(TokenType::RParen); i += 1; }
            ':' => { tokens.push(TokenType::Colon); i += 1; }
            ';' => { tokens.push(TokenType::Semicolon); i += 1; }
            ',' => { tokens.push(TokenType::Comma); i += 1; }
            '=' => { tokens.push(TokenType::Eq); i += 1; }
            '.' => { tokens.push(TokenType::Dot); i += 1; }
            '-' if i + 1 < chars.len() && chars[i + 1] == '>' => {
                tokens.push(TokenType::Arrow); i += 2;
            }
            _ if c.is_ascii_digit() => {
                let start = i;
                while i < chars.len() && chars[i].is_ascii_digit() { i += 1; }
                let num: String = chars[start..i].iter().collect();
                tokens.push(TokenType::Number(num));
            }
            _ if c.is_alphabetic() || c == '_' => {
                let start = i;
                while i < chars.len() && (chars[i].is_alphanumeric() || chars[i] == '_' || chars[i] == '<' || chars[i] == '>') { i += 1; }
                let word: String = chars[start..i].iter().collect();
                let tt = match word.as_str() {
                    "namespace" => TokenType::Namespace,
                    "struct" => TokenType::Struct,
                    "service" => TokenType::Service,
                    "rpc" => TokenType::Rpc,
                    "enum" => TokenType::Enum,
                    "stream" => TokenType::Stream,
                    _ => TokenType::Ident(word),
                };
                tokens.push(tt);
            }
            _ => { i += 1; }
        }
    }
    tokens.push(TokenType::Eof);
    tokens
}

// ─── AST ──────────────────────────────────────────────────────
#[derive(Debug)]
struct Field { tag: String, type_name: String, name: String }
#[derive(Debug)]
struct StructDef { name: String, fields: Vec<Field> }
#[derive(Debug)]
struct RpcDef { name: String, req_type: String, resp_type: String }
#[derive(Debug)]
struct ServiceDef { name: String, rpcs: Vec<RpcDef> }
#[derive(Debug)]
struct AyaFile { namespace: String, structs: Vec<StructDef>, services: Vec<ServiceDef> }

// ─── Parser ───────────────────────────────────────────────────
fn parse(tokens: &[TokenType]) -> AyaFile {
    let mut pos = 0usize;
    let mut file = AyaFile { namespace: String::new(), structs: Vec::new(), services: Vec::new() };

    while pos < tokens.len() && tokens[pos] != TokenType::Eof {
        match &tokens[pos] {
            TokenType::Namespace => {
                pos += 1;
                let mut ns = ident_val(&tokens[pos]).unwrap_or_default();
                pos += 1;
                while pos < tokens.len() && tokens[pos] == TokenType::Dot {
                    pos += 1;
                    ns.push('.');
                    ns.push_str(&ident_val(&tokens[pos]).unwrap_or_default());
                    pos += 1;
                }
                if pos < tokens.len() && tokens[pos] == TokenType::Semicolon { pos += 1; }
                file.namespace = ns;
            }
            TokenType::Enum => {
                // Skip enum blocks entirely
                pos += 1; // 'enum'
                pos += 1; // name
                if pos < tokens.len() && tokens[pos] == TokenType::LBrace {
                    pos += 1;
                    while pos < tokens.len() && tokens[pos] != TokenType::RBrace { pos += 1; }
                    if pos < tokens.len() { pos += 1; } // '}'
                }
            }
            TokenType::Struct => {
                pos += 1;
                let name = ident_val(&tokens[pos]).unwrap_or_default();
                pos += 1;
                if pos < tokens.len() && tokens[pos] == TokenType::LBrace { pos += 1; }
                let mut fields = Vec::new();
                while pos < tokens.len() && tokens[pos] != TokenType::RBrace {
                    let tag = num_val(&tokens[pos]).unwrap_or_default();
                    pos += 1;
                    if pos < tokens.len() && tokens[pos] == TokenType::Colon { pos += 1; }
                    let type_name = ident_val(&tokens[pos]).unwrap_or_default();
                    pos += 1;
                    let field_name = ident_val(&tokens[pos]).unwrap_or_default();
                    pos += 1;
                    if pos < tokens.len() && tokens[pos] == TokenType::Semicolon { pos += 1; }
                    fields.push(Field { tag, type_name, name: field_name });
                }
                if pos < tokens.len() { pos += 1; } // '}'
                file.structs.push(StructDef { name, fields });
            }
            TokenType::Service => {
                pos += 1;
                let name = ident_val(&tokens[pos]).unwrap_or_default();
                pos += 1;
                if pos < tokens.len() && tokens[pos] == TokenType::LBrace { pos += 1; }
                let mut rpcs = Vec::new();
                while pos < tokens.len() && tokens[pos] != TokenType::RBrace {
                    if tokens[pos] != TokenType::Rpc { pos += 1; continue; }
                    pos += 1; // 'rpc'
                    let rpc_name = ident_val(&tokens[pos]).unwrap_or_default();
                    pos += 1;
                    if pos < tokens.len() && tokens[pos] == TokenType::LParen { pos += 1; }
                    // skip optional 'stream' keyword
                    if pos < tokens.len() && tokens[pos] == TokenType::Stream { pos += 1; }
                    let req = ident_val(&tokens[pos]).unwrap_or_default();
                    pos += 1;
                    if pos < tokens.len() && tokens[pos] == TokenType::RParen { pos += 1; }
                    if pos < tokens.len() && tokens[pos] == TokenType::Arrow { pos += 1; }
                    if pos < tokens.len() && tokens[pos] == TokenType::LParen { pos += 1; }
                    // skip optional 'stream' keyword
                    if pos < tokens.len() && tokens[pos] == TokenType::Stream { pos += 1; }
                    let resp = ident_val(&tokens[pos]).unwrap_or_default();
                    pos += 1;
                    if pos < tokens.len() && tokens[pos] == TokenType::RParen { pos += 1; }
                    if pos < tokens.len() && tokens[pos] == TokenType::Semicolon { pos += 1; }
                    rpcs.push(RpcDef { name: rpc_name, req_type: req, resp_type: resp });
                }
                if pos < tokens.len() { pos += 1; } // '}'
                file.services.push(ServiceDef { name, rpcs });
            }
            _ => { pos += 1; }
        }
    }
    file
}

fn ident_val(t: &TokenType) -> Option<String> {
    if let TokenType::Ident(s) = t { Some(s.clone()) } else { None }
}
fn num_val(t: &TokenType) -> Option<String> {
    if let TokenType::Number(s) = t { Some(s.clone()) } else { None }
}
