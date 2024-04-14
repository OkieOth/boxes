use crate::doc::Doc;
use std::env;
use std::fs;


pub fn get_doc_from_file(file_name: &str) -> Doc {
    match fs::read_to_string(file_name) {
        Ok(s) => {
            match serde_json::from_str::<Doc>(&s) {
                Ok(d) => {
                    return d
                },
                Err(e) => {
                    panic!("Error while parsing file [{}]: {}", file_name, e);
                }
            }
        }
        Err(e) => {
            panic!("Error while reading file [{}]: {}", file_name, e);
        }
    }
}

fn get_file_name(file: &str) -> String {
    let cur_dir = env::current_dir().unwrap();
    let cur_dir_str = cur_dir.to_str().unwrap();
    format!("{}/../models/examples/{}", cur_dir_str, file)
}

#[test]
fn load_json_from_file_and_create_doc1() {
    let _ = get_doc_from_file(&get_file_name("simple_vlayout.json"));
    println!(":)");
}

#[test]
fn load_json_from_file_and_create_doc2() {
    let _ = get_doc_from_file(&get_file_name("complex_vlayout.json"));
    println!(":)");
}
