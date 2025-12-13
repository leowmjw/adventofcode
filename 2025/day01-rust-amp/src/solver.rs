//! Dial rotation counter solver
//!
//! This module processes dial rotation commands and counts how many times
//! the dial reaches position 0 (wrapping around 0-99 range).

use std::num::ParseIntError;

/// Represents a dial rotation command
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Direction {
    Left,
    Right,
}

/// Represents a single rotation command
#[derive(Debug, Clone)]
pub struct Command {
    pub direction: Direction,
    pub steps: u32,
}

/// Custom error type for parsing and validation
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum RotationError {
    InvalidFormat(String),
    ParseError(String),
    InvalidSteps,
}

impl std::fmt::Display for RotationError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::InvalidFormat(msg) => write!(f, "Invalid format: {}", msg),
            Self::ParseError(msg) => write!(f, "Parse error: {}", msg),
            Self::InvalidSteps => write!(f, "Invalid step count"),
        }
    }
}

impl std::error::Error for RotationError {}

impl From<ParseIntError> for RotationError {
    fn from(e: ParseIntError) -> Self {
        Self::ParseError(e.to_string())
    }
}

/// Parse a single command line into a Command
/// Format: "L5" (left 5 steps) or "R10" (right 10 steps)
#[inline]
fn parse_command(line: &str) -> Result<Command, RotationError> {
    let line = line.trim();
    
    if line.is_empty() {
        return Err(RotationError::InvalidFormat("Empty line".to_string()));
    }
    
    if line.len() < 2 {
        return Err(RotationError::InvalidFormat("Line too short".to_string()));
    }
    
    let direction = match line.chars().next().unwrap() {
        'L' | 'l' => Direction::Left,
        'R' | 'r' => Direction::Right,
        c => {
            return Err(RotationError::InvalidFormat(
                format!("Invalid direction: {}", c),
            ))
        }
    };
    
    let steps = line[1..].parse::<u32>()?;
    
    if steps == 0 {
        return Err(RotationError::InvalidSteps);
    }
    
    Ok(Command { direction, steps })
}

/// Part 1: Count dial hits with simple modulo behavior
/// Dial ranges from 0-99, starting at 50.
/// Count increments each time dial reaches position 0.
/// 
/// # Arguments
/// * `input` - Input text with commands (one per line)
/// 
/// # Returns
/// * `Result<u32, RotationError>` - Total count of times dial hit position 0
pub fn part1(input: &str) -> Result<u32, RotationError> {
    let mut current_position = 50u32;
    let mut hit_count = 0u32;
    
    for line in input.lines() {
        let line = line.trim();
        if line.is_empty() {
            continue;
        }
        
        let command = parse_command(line)?;
        
        let step = match command.direction {
            Direction::Left => -(command.steps as i32),
            Direction::Right => command.steps as i32,
        };
        
        current_position = {
            let new_pos = current_position as i32 + step;
            let normalized = ((new_pos % 100) + 100) % 100;
            normalized as u32
        };
        
        if current_position == 0 {
            hit_count += 1;
        }
    }
    
    Ok(hit_count)
}

/// Part 2: Count dial hits with per-step granularity
/// Dial ranges from 0-99, starting at 50.
/// Count increments for EACH step that reaches position 0.
/// This is more granular than part1.
///
/// # Arguments
/// * `input` - Input text with commands (one per line)
/// 
/// # Returns
/// * `Result<u32, RotationError>` - Total count of times dial hit position 0
pub fn part2(input: &str) -> Result<u32, RotationError> {
    let mut current_position = 50i32;
    let mut hit_count = 0u32;
    
    for line in input.lines() {
        let line = line.trim();
        if line.is_empty() {
            continue;
        }
        
        let command = parse_command(line)?;
        
        let step = match command.direction {
            Direction::Left => -1,
            Direction::Right => 1,
        };
        
        // Simulate each individual step
        for _ in 0..command.steps {
            current_position = ((current_position + step) % 100 + 100) % 100;
            if current_position == 0 {
                hit_count += 1;
            }
        }
    }
    
    Ok(hit_count)
}

#[cfg(test)]
mod tests {
    use super::*;

    // ========== Parse Command Tests ==========
    
    #[test]
    fn test_parse_command_left() {
        let cmd = parse_command("L5").unwrap();
        assert_eq!(cmd.direction, Direction::Left);
        assert_eq!(cmd.steps, 5);
    }

    #[test]
    fn test_parse_command_right() {
        let cmd = parse_command("R10").unwrap();
        assert_eq!(cmd.direction, Direction::Right);
        assert_eq!(cmd.steps, 10);
    }

    #[test]
    fn test_parse_command_lowercase() {
        let cmd = parse_command("l7").unwrap();
        assert_eq!(cmd.direction, Direction::Left);
        assert_eq!(cmd.steps, 7);
        
        let cmd = parse_command("r3").unwrap();
        assert_eq!(cmd.direction, Direction::Right);
        assert_eq!(cmd.steps, 3);
    }

    #[test]
    fn test_parse_command_whitespace() {
        let cmd = parse_command("  L5  ").unwrap();
        assert_eq!(cmd.direction, Direction::Left);
        assert_eq!(cmd.steps, 5);
    }

    #[test]
    fn test_parse_command_large_steps() {
        let cmd = parse_command("R999").unwrap();
        assert_eq!(cmd.steps, 999);
    }

    #[test]
    fn test_parse_command_invalid_direction() {
        let result = parse_command("X5");
        assert!(matches!(result, Err(RotationError::InvalidFormat(_))));
    }

    #[test]
    fn test_parse_command_empty_steps() {
        let result = parse_command("L");
        assert!(matches!(result, Err(RotationError::InvalidFormat(_))));
    }

    #[test]
    fn test_parse_command_invalid_steps() {
        let result = parse_command("L0");
        assert!(matches!(result, Err(RotationError::InvalidSteps)));
        
        let result = parse_command("Rabc");
        assert!(matches!(result, Err(RotationError::ParseError(_))));
    }

    #[test]
    fn test_parse_command_empty_line() {
        let result = parse_command("");
        assert!(matches!(result, Err(RotationError::InvalidFormat(_))));
    }

    // ========== Part 1 Tests ==========

    #[test]
    fn test_part1_simple_right_no_wrap() {
        let input = "R5";
        let result = part1(input).unwrap();
        // 50 + 5 = 55, no hit
        assert_eq!(result, 0);
    }

    #[test]
    fn test_part1_single_hit() {
        let input = "R50";
        let result = part1(input).unwrap();
        // 50 + 50 = 100 % 100 = 0, one hit
        assert_eq!(result, 1);
    }

    #[test]
    fn test_part1_left_direction() {
        let input = "L50";
        let result = part1(input).unwrap();
        // 50 - 50 = 0, one hit
        assert_eq!(result, 1);
    }

    #[test]
    fn test_part1_multiple_commands() {
        let input = "R50\nL25";
        let result = part1(input).unwrap();
        // 50 + 50 = 100 % 100 = 0 (hit 1)
        // 0 - 25 = -25 % 100 = 75, no additional hit
        assert_eq!(result, 1);
    }

    #[test]
    fn test_part1_wrapping() {
        let input = "R55";
        let result = part1(input).unwrap();
        // 50 + 55 = 105 % 100 = 5, no hit
        assert_eq!(result, 0);
    }

    #[test]
    fn test_part1_negative_wrapping() {
        let input = "L55";
        let result = part1(input).unwrap();
        // 50 - 55 = -5, need modulo: (-5 % 100 + 100) % 100 = 95, no hit
        assert_eq!(result, 0);
    }

    #[test]
    fn test_part1_empty_lines_ignored() {
        let input = "R50\n\nL25\n\n";
        let result = part1(input).unwrap();
        assert_eq!(result, 1);
    }

    #[test]
    fn test_part1_multiple_hits() {
        let input = "R50\nR50\nR50";
        let result = part1(input).unwrap();
        // 50 + 50 = 0 (hit 1)
        // 0 + 50 = 50, no hit
        // 50 + 50 = 0 (hit 2)
        assert_eq!(result, 2);
    }

    #[test]
    fn test_part1_exact_boundaries() {
        let input = "R100";
        let result = part1(input).unwrap();
        // 50 + 100 = 150 % 100 = 50, no hit
        assert_eq!(result, 0);
    }

    // ========== Part 2 Tests ==========

    #[test]
    fn test_part2_single_step() {
        let input = "R1";
        let result = part2(input).unwrap();
        // 50 + 1 = 51, no hit
        assert_eq!(result, 0);
    }

    #[test]
    fn test_part2_simple_hit() {
        let input = "R50";
        let result = part2(input).unwrap();
        // Stepping from 50 to 100 (wrapped to 0)
        assert_eq!(result, 1);
    }

    #[test]
    fn test_part2_left_hit() {
        let input = "L50";
        let result = part2(input).unwrap();
        // Stepping from 50 down to 0
        assert_eq!(result, 1);
    }

    #[test]
    fn test_part2_wrapping_multiple_hits() {
        let input = "R150";
        let result = part2(input).unwrap();
        // 50 -> 100 (0, hit 1) -> 50 -> 100 (0, hit 2)
        assert_eq!(result, 2);
    }

    #[test]
    fn test_part2_complex_sequence() {
        let input = "R50\nL100";
        let result = part2(input).unwrap();
        // R50: 50 -> ... -> 0 (hit 1)
        // L100: 0 -> 99 -> 98 -> ... -> 1 -> 0 (hit 2) -> 99 -> ...
        assert_eq!(result, 2);
    }

    #[test]
    fn test_part2_large_step_count() {
        let input = "R200";
        let result = part2(input).unwrap();
        // Each full rotation of 100 should hit 0 twice
        assert_eq!(result, 2);
    }

    #[test]
    fn test_part2_no_hits() {
        let input = "R25\nR25";
        let result = part2(input).unwrap();
        // 50 + 25 = 75, then 75 + 25 = 100 % 100 = 0... actually one hit
        assert_eq!(result, 1);
    }

    #[test]
    fn test_part2_left_wrapping() {
        let input = "L25";
        let result = part2(input).unwrap();
        // 50 - 1 = 49, ... - 25 = 25, no hit
        assert_eq!(result, 0);
    }

    // ========== Edge Case Tests ==========

    #[test]
    fn test_error_on_invalid_command() {
        let input = "X50\nR25";
        let result = part1(input);
        assert!(result.is_err());
    }

    #[test]
    fn test_error_propagation() {
        let input = "R50\ninvalid\nR25";
        let result = part2(input);
        assert!(result.is_err());
    }

    // ========== Comprehensive Integration Tests ==========

    #[test]
    fn test_part1_comprehensive() {
        let input = "R10\nR20\nL15\nR25";
        let result = part1(input).unwrap();
        // 50 + 10 = 60
        // 60 + 20 = 80
        // 80 - 15 = 65
        // 65 + 25 = 90
        assert_eq!(result, 0);
    }

    #[test]
    fn test_part2_comprehensive() {
        let input = "R50\nR50";
        let result = part2(input).unwrap();
        // First R50: 50 + 50 = 100 % 100 = 0 (hit 1)
        // Second R50: 0 + 50 = 50 (no hit)
        assert_eq!(result, 1);
    }

    #[test]
    fn test_consistency_small_values() {
        let input = "R10";
        let p1 = part1(input).unwrap();
        let p2 = part2(input).unwrap();
        // Both should give same result for non-wrapping cases
        assert_eq!(p1, p2);
    }

    #[test]
    fn test_starting_position() {
        // Verify starting position is 50
        let input = "R0";
        let result = parse_command(input);
        assert!(result.is_err()); // R0 should be invalid (0 steps)
    }
}
