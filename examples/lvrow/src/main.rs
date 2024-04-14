use std::env;

use bimpl::doc::Doc;
use bimpl::serialize::get_doc_from_file;


// this is an comment
fn main() {
    let cur_dir = env::current_dir().unwrap();
    let cur_dir_str = cur_dir.to_str().unwrap();

    let file_name =format!("{}/models/examples/complex_vlayout.json", cur_dir_str);

    let doc: Doc = get_doc_from_file(&file_name);
    println!("{:?}", doc);
}

#[cfg(test)]
mod tests {
    use bimpl::traits::Content;
    use builder_m4cro::Builder;
    
    #[derive(Builder, Debug)]
    pub struct MyContent {
        caption: String,
        id: usize,
        text1: String,
        text2: String,
    }
    
    impl Content for MyContent {
        fn caption(&self) -> &String {
            &self.caption
        }
        fn id(&self) -> usize {
            self.id
        }
        fn text1(&self) -> &String {
            &self.text1
        }
        fn text2(&self) -> &String {
            &self.text2
        }
    }
    

    #[test]
    fn test_builder() {
        let c1: MyContent = MyContent::builder()
            .caption("Caption1".to_string())
            .id(1)
            .text1("This is a longer text".to_string())
            .text2("This is a second text, even longer then the first one".to_string())
            .build()
            .unwrap();

        let c2: Result<MyContent, String> = MyContent::builder()
            .caption("My Caption2".to_string())
            .build();
        let t = c1.text1();
        println!("Hello-1: {:?}", c1);

        println!("Hello-2: {:?}", c2);
    }
}
