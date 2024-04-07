
pub trait Content {
    fn caption(&self) -> &String;
    fn id(&self) -> usize;
    fn text1(&self) -> &String;
    fn text2(&self) -> &String;
}