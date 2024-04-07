use std::fmt;
use std::fmt::{Debug, Formatter};

use crate::traits::Content;

pub struct SimpleCell {
    pub x: usize,
    pub y: usize,
    pub width: usize,
    pub height: usize,
    pub content: Box<dyn Content>,
}

impl Debug for SimpleCell {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        f.debug_struct("SimpleCell")
            .field("width", &self.width)
            .field("height", &self.height)
            // You might want to customize how you print content depending on its type
            .field(
                "content",
                &format_args!("{}:{}", &self.content.id(), &self.content.caption()),
            )
            .finish()
    }
}

#[derive(Debug)]
pub struct Horizontal {
    pub parts: Vec<Box<Layout>>,
}

#[derive(Debug)]
pub struct Vertical {
    pub parts: Vec<Box<Layout>>,
}

#[derive(Debug)]
pub struct Grid {
    pub parts: Vec<Vec<Box<Layout>>>,
}

#[derive(Debug)]
pub enum Layout {
    Simple(SimpleCell),
    Vertical(Vertical),
    Horizontal(Horizontal),
    Grid(Grid),
}
