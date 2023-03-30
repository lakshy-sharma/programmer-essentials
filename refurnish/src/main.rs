/*
Copyright © [2022] [Lakshy Sharma] <lakshy.sharma@protonmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
use std::env;
use std::fs;
use text_colorizer::*;
use regex:: Regex;

// Declaring the structure required for storing input arguements.
#[derive(Debug)]
struct Arguements {
    target_pattern: String,
    replacement_pattern: String,
    input_file: String,
    output_file: String,
}

fn print_usage() {
    // Printing a message for invalid command invoking.
    let util_name = "refurnish".blue().bold();
    eprintln!("{util_name} - Change one string pattern into another for any file!");
    eprintln!("Usage: refurnish <target_pattern> <replacement_pattern> <input_filename> <output_filename>");
}

fn parse_args() -> Arguements {
    // A function to read all the arguements and validate their availablity.

    // Collect the input arguements from the cli and ensure they are all stored in a vector containing strings.
    // Using .skip() function to avoid reading the program name as starting variable and .collect() to iterate over many arguements.
    let args: Vec<String> = env::args().skip(1).collect();

    // Throw an exception if the arguements provided are less than required.
    if args.len() != 4 {
        print_usage();
        eprintln!("{}: wrong number of arguements, expected 4 but got {} ", "Error".red().bold(), args.len());
        std::process::exit(1)
    }

    // Copy arguement values in the structure if the arguements are available.
    Arguements {
        target_pattern: args[0].clone(),
        replacement_pattern: args[1].clone(),
        input_file: args[2].clone(),
        output_file: args[3].clone()
    }
}

fn replace(target_pattern: &str, replacement_pattern: &str, text: &str) -> Result<String, regex::Error>{
    // Declaring our new regex pattern.
    let regex=Regex::new(target_pattern)?;
    Ok(regex.replace_all(text,replacement_pattern).to_string())

}

fn main() {
    // Capturing the arguements.
    let args = parse_args();
    
    // Reading the input data from file in string format.
    let input_data = match fs::read_to_string(&args.input_file) {
        // If the function returns a string then it is fine.
        Ok(v) => v,
        // If the function returns a erro then handle it.
        Err(e) => {
            eprintln!("{}: Failed to read from the file {}.\n{}: {}", "Error".red().bold(),args.input_file, "Message".green(), e);
            std::process::exit(1);
        }
    };

    // Replacing the patterns in the data by using regex replacer function we cretaed.
    let edited_data = match replace(&args.target_pattern, &args.replacement_pattern, &input_data) {
        Ok(v) => v,
        Err(e) => {
            eprintln!("{}: Failed to replace the data th text: {:?}", "Error".red().bold(), e);
            std::process::exit(1);
        }
    };

    // Writing the data into a output file.
    match fs::write(&args.output_file, &edited_data) {
        Ok(_) => {},
        Err(e) => {
            eprintln!("{}: Failed to write to the file {}.\n{}: {}", "Error".red().bold(),args.output_file,"Message".green(), e);
            std::process::exit(1);
        }
    }
}