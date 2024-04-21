use std::default;

use crate::layouts;
use serde::{Deserialize, Serialize};

use crate::text_helper::FontGroup;


#[derive(Deserialize, Serialize, Debug, Default)]
pub enum TextAlignment {
    Left,
    Right,
    #[default]
    Center,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct TextConfig {
    pub margin_top: u32,
    pub size: u32,
    pub font: FontGroup,
    pub alignment: TextAlignment,
}

impl Default for TextConfig {
    fn default() -> Self {
        TextConfig {
            margin_top: 10,
            size: 12,
            font: FontGroup::Sans,
            alignment: TextAlignment::default(),
        }
    }
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

impl Default for DocConfig {
    fn default() -> Self {
        DocConfig {
            default_max_box_with: 200,
            text_margin: 10,
            line_space: 1.5,
            font_caption: TextConfig::default(),
            font_text1: TextConfig::default(),
            font_text2: TextConfig::default(),
        }
    }
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Doc {
    pub layout: layouts::Layout,
    pub config: Option<DocConfig>
}
