use crate::doc::Doc;
use std::fmt;
use std::env;
use std::fs;
use std::path::Path;
use std::path::PathBuf;

#[test]
fn load_json_from_file_and_create_doc() {
    // open a json file
    // if let Ok(pb) = env::current_dir() {
    //     pb.
    // } else {
    //     panic!("can't query current dir");
    // }
    let cur_dir = env::current_dir().unwrap();
    let cur_dir_str = cur_dir.to_str().unwrap();
    let file_name = format!("{}/../models/examples/simple_vlayout.json", cur_dir_str);
    match fs::read_to_string(&file_name) {
        Ok(s) => {
            match serde_json::from_str::<Doc>(&s) {
                Ok(d) => {
                    println!(":)")
                },
                Err(e) => {
                    panic!("Error while parsing file [{}]: {}", &file_name, e);
                }
            }
        }
        Err(e) => {
            panic!("Error while reading file [{}]: {}", &file_name, e);
        }
    }
}
