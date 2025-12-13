// use std::io::BufRead;

fn main() {

    println!("Day01 - Rust");
    // part1();
    let my_str = include_str!("part1.txt");
    part2(my_str);
}


// The output is wrapped in a Result to allow matching on errors.
// Returns an Iterator to the Reader of the lines of the file.
// fn read_lines<P>(filename: P) -> io::Result<io::Lines<io::BufReader<File>>>
// where P: AsRef<Path>, {
//     let file = File::open(filename)?;
//     Ok(io::BufReader::new(file).lines())
// }

fn part2(my_str: &str) {
    // Below will be the more complex with steps more thsn 100++
    // let my_str = include_str!("part2a.txt");

    let mut count = 0;
    let mut current = 50;

    my_str.lines().for_each(
        |line| {
            // Extratc first char actiion .. maybe use the Match?
            let cmd = &line[0..1];
            // dbg!(cmd);
            // Extract out the steps as int ..
            let steps = line[1..].parse().unwrap();
            let steps_size;
            // If action is 'L'; then compare if eq, less, more ..
            if cmd == "L" {
                steps_size = -1;
            } else {
                steps_size = 1;
            }

            // For number of steps; simukate throuh each
            // Once it recahes to 0 when modules then increase the counter
            for _d in 1..=steps {
                // dbg!(current);
                current = ((current + steps_size) % 100 + 100 ) % 100;
                if current == 0 {
                    count = count + 1;
                    println!("HIT!! {}", count);
                }
                // DEBUG
                // println!("CURRENT: {} AFTER STEP: {}", current, steps_size);
            };
        }
    );
    println!("COUNT: {}", count);
}

fn part1() {

    // let myinput = "test.txt";
    // let my_str = include_str!("test.txt");

    let my_str = include_str!("part1.txt");

    // println!("my_str: {}", my_str);

    // include_str!("test.txt");
    // Open file and filter line by line ...
    // File::open(input).lines().filter(|line| println!("{}", line));

    let mut count = 0;
    let mut current = 50;

    my_str.lines().for_each(
        |line| {
            // Extratc first char actiion .. maybe use the Match?
            let cmd = &line[0..1];
            // Extract out the steps as int ..
            let steps = &line[1..].parse::<i32>().unwrap();
            // If action is 'L'; then compare if eq, less, more ..
            if cmd == "L" {
                // DEBUG
                current = current + (steps * -1);
                println!("ACT: {}", steps * -1);
                println!("CURRENT: {}", current);
            } else {
                current = current + steps;
                println!("ACT: {}", steps);
                println!("CURRENT: {}", current);
            }
            // Rule is simpler .. no need modulus
            // If it reaches exactly 0; the increase the counter
            if current % 100 == 0 {
                count = count + 1;
                println!("HIT!! {}", count);
            }
        }
    );

    println!("COUNT: {}", count);


    // read_lines(input).unwrap().for_each(|line| println!("{}", line.unwrap()));

    // Extract out each line and print it out


}
