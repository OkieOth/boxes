use serde::Deserialize;
use std::fmt;
use std::fmt::{Debug, Formatter};


#[derive(Debug, Deserialize)]
pub struct Content {
    pub id: String,
    pub caption: String,
    pub text1: Option<String>,
    pub text2: Option<String>,
}

#[derive(Debug, Deserialize)]
pub struct Geometry {
    pub x: usize,
    pub y: usize,
    pub width: usize,
    pub height: usize,
}

#[derive(Debug, Deserialize)]
pub struct SimpleCell {
    pub geometry: Option<Geometry>,
    pub content: Content,
}

#[derive(Debug, Deserialize)]
pub enum Layout {
    #[serde(rename(deserialize = "simple"))]
    Simple(SimpleCell),

    #[serde(rename(deserialize = "vertical"))]
    Vertical(Vec<Layout>),

    #[serde(rename(deserialize = "horizontal"))]
    Horizontal(Vec<Layout>),

    #[serde(rename(deserialize = "grid"))]
    Grid(Vec<Vec<Layout>>),
}
