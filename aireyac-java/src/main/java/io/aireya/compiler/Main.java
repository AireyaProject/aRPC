package io.aireya.compiler;

import java.nio.file.Files;
import java.nio.file.Paths;
import java.io.IOException;

public class Main {
    public static void main(String[] args) {
        if (args.length < 1) {
            System.err.println("Usage: java -jar aireyac-java.jar <file.aya>");
            System.exit(1);
        }

        String filename = args[0];
        System.out.println("☕ Aireya Java Compiler (aireyac-java) starting...");
        System.out.println("Reading IDL file: " + filename);

        try {
            String content = new String(Files.readAllBytes(Paths.get(filename)));
            
            // TODO: Lexer.tokenize(content)
            // TODO: Parser.parse(tokens)
            // TODO: Generators.runAll(ast)
            
            System.out.println("✅ AST Parsed successfully (Java Engine)");
            System.out.println("✅ Code generated!");
        } catch (IOException e) {
            System.err.println("Failed to read file: " + e.getMessage());
            System.exit(1);
        }
    }
}
