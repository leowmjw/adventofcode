mod solver;

use solver::{part1, part2};

fn main() {
    println!("Day 01 - Dial Rotation Counter");

    // Example input
    let test_input = include_str!("../input.txt");
    
    println!("\n=== Part 1 ===");
    match part1(test_input) {
        Ok(result) => println!("Result: {}", result),
        Err(e) => eprintln!("Error: {}", e),
    }
    
    println!("\n=== Part 2 ===");
    match part2(test_input) {
        Ok(result) => println!("Result: {}", result),
        Err(e) => eprintln!("Error: {}", e),
    }
}
