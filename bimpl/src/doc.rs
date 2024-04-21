use crate::layouts;
use serde::{Deserialize, Serialize};


#[derive(Deserialize, Serialize, Debug)]
pub enum TextAlignment {
    Left,
    Right,
    Center,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct TextConfig {
    pub margin_top: u32,
    pub size: u32,
    pub font: String,
    pub alignment: TextAlignment,
}


#[derive(Deserialize, Serialize, Debug)]
pub struct DocConfig {
    /// Initial max width of boxes. After processing a box can be longer.
    pub default_max_box_with: u32,

    /// Margin of text boxes to the outher box border
    pub text_margin: u32,

    /// Space between lines in one text block in relation to the font size, e.g. 1.5
    pub line_space: f32,

    /// text config for the general caption
    pub font_caption: TextConfig,

    /// text config for text 1
    pub font_text1: TextConfig,

    /// text config for text 2
    pub font_text2: TextConfig,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Doc {
    pub layout: layouts::Layout,
    pub config: Option<DocConfig>
}
