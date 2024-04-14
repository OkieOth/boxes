use rusttype::{point, Font, Scale};

pub enum FontGroup {
    Sans,
    SansItalic,
    SansBold,
    SansBoldItalic,
    Mono,
    MonoItalic,
}

pub fn get_text_bbox(text: &str, font_size: usize, font_group: FontGroup) -> (usize, usize) {
    
}