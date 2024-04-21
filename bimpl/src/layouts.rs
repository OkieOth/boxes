use serde::{Deserialize, Serialize};
use std::fmt::Debug;

use crate::doc::DocConfig;


#[derive(Debug, Deserialize, Serialize)]
pub struct Geometry {
    pub x: usize,
    pub y: usize,
    pub width: usize,
    pub height: usize,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct SimpleCell {
    pub geometry: Option<Geometry>,
    pub id: String,
    pub caption: String,
    pub text1: Option<String>,
    pub text2: Option<String>,
}


#[derive(Debug, Deserialize, Serialize)]
pub enum Layout {
    #[serde(rename = "simple")]
    Simple(SimpleCell),

    #[serde(rename = "vertical")]
    Vertical(Vec<Layout>),

    #[serde(rename = "horizontal")]
    Horizontal(Vec<Layout>),

    #[serde(rename = "grid")]
    Grid(Vec<Vec<Layout>>),
}

#[cfg(test)]
mod tests {
    use std::io::Error;

    use super::*;
    use crate::doc::Doc;

    #[test]
    fn deserialize_test() {
        let s1 = Layout::Simple(SimpleCell {
            geometry: None,
            id: "id1".to_string(),
            caption: "Box 1".to_string(),
            text1: Some("This is an additional text".to_string()),
            text2: Some("Here is a second text, that is a little bit larger than the first one".to_string()),
        });
        let s2 = Layout::Simple(SimpleCell {
            geometry: None,
            id: "id2".to_string(),
            caption: "Box 1".to_string(),
            text1: Some("comment v2".to_string()),
            text2: None,
        });
        let s3 = Layout::Simple(SimpleCell {
            geometry: None,
            id: "id3".to_string(),
            caption: "Box 1".to_string(),
            text1: None,
            text2: None,
        });
        let vert = Layout::Vertical(vec![s1, s2, s3]);
        let d = Doc {
            layout: vert,
            config: None,
        };

        let json_string = serde_json::to_string(&d).unwrap();
        write_to_file(&json_string, "test_output.json").unwrap();
    }

    fn write_to_file(json_string: &str, file_name: &str) -> Result<(), Error>{
        use std::fs::File;
        use std::io::Write;
        use std::env;


        let cur_dir = env::current_dir().unwrap();
        let cur_dir_str = cur_dir.to_str().unwrap();
        // "/home/eiko/prog/git/boxes/examples/lvrow"
    
        let file_name = if cur_dir_str.ends_with("bimpl") {
            format!("{}/../tmp/{}", cur_dir_str, file_name)
        } else {
            format!("{}/tmp/{}", cur_dir_str, file_name)
        };
    
        let mut file = File::create(file_name)?;
        file.write_all(json_string.as_bytes())?;
        Ok(())
    }
}
