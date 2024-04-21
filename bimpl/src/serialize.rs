use crate::doc::Doc;
use std::fs;

pub fn get_doc_from_file(file_name: &str) -> Doc {
    match fs::read_to_string(file_name) {
        Ok(s) => match serde_json::from_str::<Doc>(&s) {
            Ok(d) => return d,
            Err(e) => {
                panic!("Error while parsing file [{}]: {}", file_name, e);
            }
        },
        Err(e) => {
            panic!("Error while reading file [{}]: {}", file_name, e);
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::env;
    use crate::layouts::Layout;

    fn get_file_name(file: &str) -> String {
        let cur_dir = env::current_dir().unwrap();
        let cur_dir_str = cur_dir.to_str().unwrap();
        format!("{}/../models/examples/{}", cur_dir_str, file)
    }

    #[test]
    fn load_json_from_file_and_create_doc1() {
        let doc = get_doc_from_file(&get_file_name("simple_vlayout.json"));
        let l = doc.layout;
        match l {
            Layout::Vertical(v) => {
                v.iter().for_each(|cl| {
                    if let Layout::Simple(s) = cl {
                        assert!(s.id != "");
                        assert!(s.caption != "");
                    } else {
                        panic!("expected simple layout, got something different");
                    }
                }) 
            },
            Layout::Horizontal(_h) => {
                panic!("expected vertical layout, got horizontal layout");
            },
            Layout::Grid(_g) => {
                panic!("expected vertical layout, got grid layout");
            },
            Layout::Simple(_s) => {
                panic!("expected vertical layout, got simple layout");
            }
        }
    }

    #[test]
    fn load_json_from_file_and_create_doc2() {
        let _ = get_doc_from_file(&get_file_name("complex_vlayout.json"));
        println!(":)");
    }
}
