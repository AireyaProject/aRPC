package io.aireya.compiler;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

public class Main {
    // ─── Token ────────────────────────────────────────────────
    enum TokenType { NAMESPACE, STRUCT, SERVICE, RPC, ENUM, STREAM,
        LBRACE, RBRACE, LPAREN, RPAREN, COLON, SEMICOLON, ARROW, COMMA, EQ, DOT,
        IDENT, NUMBER, EOF }

    record Token(TokenType type, String value) {}

    // ─── Lexer ────────────────────────────────────────────────
    static List<Token> tokenize(String src) {
        List<Token> tokens = new ArrayList<>();
        int i = 0;
        while (i < src.length()) {
            char c = src.charAt(i);
            if (Character.isWhitespace(c)) { i++; continue; }
            if (c == '/' && i + 1 < src.length() && src.charAt(i + 1) == '/') {
                while (i < src.length() && src.charAt(i) != '\n') i++;
                continue;
            }
            if (c == '{') { tokens.add(new Token(TokenType.LBRACE, "{")); i++; continue; }
            if (c == '}') { tokens.add(new Token(TokenType.RBRACE, "}")); i++; continue; }
            if (c == '(') { tokens.add(new Token(TokenType.LPAREN, "(")); i++; continue; }
            if (c == ')') { tokens.add(new Token(TokenType.RPAREN, ")")); i++; continue; }
            if (c == ':') { tokens.add(new Token(TokenType.COLON, ":")); i++; continue; }
            if (c == ';') { tokens.add(new Token(TokenType.SEMICOLON, ";")); i++; continue; }
            if (c == ',') { tokens.add(new Token(TokenType.COMMA, ",")); i++; continue; }
            if (c == '=') { tokens.add(new Token(TokenType.EQ, "=")); i++; continue; }
            if (c == '.') { tokens.add(new Token(TokenType.DOT, ".")); i++; continue; }
            if (c == '-' && i + 1 < src.length() && src.charAt(i + 1) == '>') {
                tokens.add(new Token(TokenType.ARROW, "->"));
                i += 2; continue;
            }
            if (Character.isDigit(c)) {
                int start = i;
                while (i < src.length() && Character.isDigit(src.charAt(i))) i++;
                tokens.add(new Token(TokenType.NUMBER, src.substring(start, i)));
                continue;
            }
            if (Character.isLetter(c) || c == '_') {
                int start = i;
                while (i < src.length() && (Character.isLetterOrDigit(src.charAt(i)) || src.charAt(i) == '_' || src.charAt(i) == '<' || src.charAt(i) == '>')) i++;
                String word = src.substring(start, i);
                TokenType tt = switch (word) {
                    case "namespace" -> TokenType.NAMESPACE;
                    case "struct" -> TokenType.STRUCT;
                    case "service" -> TokenType.SERVICE;
                    case "rpc" -> TokenType.RPC;
                    case "enum" -> TokenType.ENUM;
                    case "stream" -> TokenType.STREAM;
                    default -> TokenType.IDENT;
                };
                tokens.add(new Token(tt, word));
                continue;
            }
            i++;
        }
        tokens.add(new Token(TokenType.EOF, ""));
        return tokens;
    }

    // ─── AST ──────────────────────────────────────────────────
    record Field(String tag, String type, String name) {}
    record StructDef(String name, List<Field> fields) {}
    record RpcDef(String name, String reqType, String respType) {}
    record ServiceDef(String name, List<RpcDef> rpcs) {}

    // ─── Parser ───────────────────────────────────────────────
    static int pos = 0;

    static Token peek(List<Token> tokens) { return tokens.get(pos); }
    static Token advance(List<Token> tokens) { return tokens.get(pos++); }
    static void expect(List<Token> tokens, TokenType t) {
        if (advance(tokens).type() != t) throw new RuntimeException("Unexpected token at pos " + (pos - 1));
    }

    static String parseNamespace(List<Token> tokens) {
        advance(tokens); // 'namespace'
        StringBuilder ns = new StringBuilder(advance(tokens).value());
        while (peek(tokens).type() == TokenType.DOT) {
            advance(tokens);
            ns.append(".").append(advance(tokens).value());
        }
        expect(tokens, TokenType.SEMICOLON);
        return ns.toString();
    }

    static StructDef parseStruct(List<Token> tokens) {
        advance(tokens); // 'struct'
        String name = advance(tokens).value();
        expect(tokens, TokenType.LBRACE);
        List<Field> fields = new ArrayList<>();
        while (peek(tokens).type() != TokenType.RBRACE) {
            String tag = advance(tokens).value();
            expect(tokens, TokenType.COLON);
            String type = advance(tokens).value();
            String fieldName = advance(tokens).value();
            expect(tokens, TokenType.SEMICOLON);
            fields.add(new Field(tag, type, fieldName));
        }
        advance(tokens); // '}'
        return new StructDef(name, fields);
    }

    static ServiceDef parseService(List<Token> tokens) {
        advance(tokens); // 'service'
        String name = advance(tokens).value();
        expect(tokens, TokenType.LBRACE);
        List<RpcDef> rpcs = new ArrayList<>();
        while (peek(tokens).type() != TokenType.RBRACE) {
            advance(tokens); // 'rpc'
            String rpcName = advance(tokens).value();
            expect(tokens, TokenType.LPAREN);
            String reqType = advance(tokens).value();
            expect(tokens, TokenType.RPAREN);
            expect(tokens, TokenType.ARROW);
            expect(tokens, TokenType.LPAREN);
            String respType = advance(tokens).value();
            expect(tokens, TokenType.RPAREN);
            expect(tokens, TokenType.SEMICOLON);
            rpcs.add(new RpcDef(rpcName, reqType, respType));
        }
        advance(tokens); // '}'
        return new ServiceDef(name, rpcs);
    }

    // ─── Main ─────────────────────────────────────────────────
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
            List<Token> tokens = tokenize(content);

            List<StructDef> structs = new ArrayList<>();
            List<ServiceDef> services = new ArrayList<>();
            String namespace = "";

            pos = 0;
            while (peek(tokens).type() != TokenType.EOF) {
                switch (peek(tokens).type()) {
                    case NAMESPACE -> namespace = parseNamespace(tokens);
                    case STRUCT -> structs.add(parseStruct(tokens));
                    case SERVICE -> services.add(parseService(tokens));
                    default -> advance(tokens);
                }
            }

            System.out.println("✅ AST Parsed: " + structs.size() + " structs, " + services.size() + " services");
            System.out.println("   Namespace: " + namespace);
            for (var s : structs) {
                System.out.println("   Struct: " + s.name() + " (" + s.fields().size() + " fields)");
            }
            for (var s : services) {
                System.out.println("   Service: " + s.name() + " (" + s.rpcs().size() + " RPCs)");
            }
            System.out.println("✅ Compilation complete.");
        } catch (IOException e) {
            System.err.println("Failed to read file: " + e.getMessage());
            System.exit(1);
        }
    }
}
