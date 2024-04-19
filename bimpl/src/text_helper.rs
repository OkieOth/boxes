use std::env;
use std::fs;
use rusttype::{point, Font, Scale};

pub enum FontGroup {
    Sans,
    SansItalic,
    SansBold,
    SansBoldItalic,
    Mono,
    MonoItalic,
}

/// Splits a given text in lines. Returns a vector of lines.
/// ```rust
/// use bimpl::text_helper::{split_text_in_lines, FontGroup};
/// 
/// let lines = split_text_in_lines("hallo ballo", 20, 18,FontGroup::Mono);
/// assert_eq!(lines.len(), 2);
/// ```
pub fn split_text_in_lines(text: &str, max_with: u32, font_size: usize, font_group: FontGroup) -> Vec<String> {
    split_text_in_lines_by_ref(text, max_with, font_size, &font_group)
}

/// Splits a given text in lines. Returns a vector of lines.
/// ```rust
/// use bimpl::text_helper::{split_text_in_lines_by_ref, FontGroup};
/// 
/// let fg = FontGroup::Mono;
/// let lines = split_text_in_lines_by_ref("hallo ballo", 20, 18,&fg);
/// assert_eq!(lines.len(), 2);
/// ```
pub fn split_text_in_lines_by_ref(text: &str, max_with: u32, font_size: usize, font_group: &FontGroup) -> Vec<String> {
    let mut lines = Vec::new();
    let mut line = String::new(); // String::new();
    for word in text.split_whitespace() {
        let mut tmp_str = line.clone();
        if tmp_str.len() > 0 {
            line.push_str(" ");
        };
        tmp_str.push_str(word);
        let with = get_text_bbox_by_ref(&tmp_str, font_size, &font_group);
        if with.0 > max_with {
            if line.len() > 0 {
                lines.push(line);
            } else {
                lines.push(tmp_str);
            }
            line = String::new();
        } else {
            line = tmp_str;
        }
    }
    if line.len() > 0 {
        lines.push(line);
    }
    lines
}


/// Returns (width, height) of a given text
/// ```rust
/// use bimpl::text_helper::{get_text_bbox_by_ref, FontGroup};
/// 
/// let fg = FontGroup::Mono;
/// let (w1,h1) = get_text_bbox_by_ref("hallo", 18,&fg);
/// assert_eq!((w1, h1),(53, 18));
/// ```
pub fn get_text_bbox(text: &str, font_size: usize, font_group: FontGroup) -> (u32, u32) {
    get_text_bbox_by_ref(text, font_size, &font_group)
}

/// Returns (width, height) of a given text
/// ```rust
/// use bimpl::text_helper::{get_text_bbox, FontGroup};
/// 
/// let (w1,h1) = get_text_bbox("hallo", 18,FontGroup::Mono);
/// assert_eq!((w1, h1),(53, 18));
/// ```
pub fn get_text_bbox_by_ref(text: &str, font_size: usize, font_group: &FontGroup) -> (u32, u32) {
    // Example from here: https://gitlab.redox-os.org/redox-os/rusttype/-/blob/master/dev/examples/image.rs

    let scale = Scale::uniform(font_size as f32);

    let font_file_bytes = get_font_file(font_group);
    let font = Font::try_from_bytes(&font_file_bytes).unwrap();

    let v_metrics = font.v_metrics(scale);
    let glyphs: Vec<_> = font.layout(text, scale, point(0.0, 0.0)).collect();
    let glyphs_height = (v_metrics.ascent - v_metrics.descent).ceil() as u32;
    let glyphs_width = {
        let min_x = glyphs
            .first()
            .map(|g| g.pixel_bounding_box().unwrap().min.x)
            .unwrap();
        let max_x = glyphs
            .last()
            .map(|g| g.pixel_bounding_box().unwrap().max.x)
            .unwrap();
        (max_x - min_x) as u32
    };
    (glyphs_width, glyphs_height)
}

fn get_font_file(font_group: &FontGroup) -> Vec<u8> {
    let cur_dir = env::current_dir().unwrap();
    let cur_dir_str = cur_dir.to_str().unwrap();

    let font_file = match *font_group {
        FontGroup::Sans => "FreeSans.ttf",
        FontGroup::SansItalic => "FreeSansOblique.ttf",
        FontGroup::SansBold => "FreeSansBold.ttf",
        FontGroup::SansBoldItalic => "FreeSansBoldOblique.ttf",
        FontGroup::Mono => "FreeMono.ttf",
        FontGroup::MonoItalic => "FreeMonoOblique.ttf",
    };

    let s = format!("{}/resources/fonts/{}", cur_dir_str, font_file);

    fs::read(s).expect("Failed to read font file")
}

#[test]
fn test_get_text_bbox_mono() {
    let (w1,h1) = get_text_bbox("hallo", 18,FontGroup::Mono);
    let (w2,h2) = get_text_bbox("hallo", 24,FontGroup::Mono);
    let (w3,h3) = get_text_bbox("hallo", 36,FontGroup::Mono);
    assert_eq!((w1, h1),(53, 18));
    assert_eq!((w2, h2),(70, 24));
    assert_eq!((w3, h3),(105, 36));
}

#[test]
fn test_get_text_bbox_mono_italic() {
    let (w1,h1) = get_text_bbox("hallo", 18,FontGroup::MonoItalic);
    let (w2,h2) = get_text_bbox("hallo", 24,FontGroup::MonoItalic);
    let (w3,h3) = get_text_bbox("hallo", 36,FontGroup::MonoItalic);
    assert_eq!((w1, h1),(54, 18));
    assert_eq!((w2, h2),(71, 24));
    assert_eq!((w3, h3),(107, 36));
}

#[test]
fn test_get_text_bbox_sans() {
    let (w1,h1) = get_text_bbox("hallo", 18,FontGroup::Sans);
    let (w2,h2) = get_text_bbox("hallo", 24,FontGroup::Sans);
    let (w3,h3) = get_text_bbox("hallo", 36,FontGroup::Sans);
    assert_eq!((w1, h1),(36, 18));
    assert_eq!((w2, h2),(48, 24));
    assert_eq!((w3, h3),(71, 36));
}

#[test]
fn test_get_text_bbox_sans_italic() {
    let (w1,h1) = get_text_bbox("hallo", 18,FontGroup::SansItalic);
    let (w2,h2) = get_text_bbox("hallo", 24,FontGroup::SansItalic);
    let (w3,h3) = get_text_bbox("hallo", 36,FontGroup::SansItalic);
    assert_eq!((w1, h1),(38, 18));
    assert_eq!((w2, h2),(50, 24));
    assert_eq!((w3, h3),(75, 36));
}

#[test]
fn test_get_text_bbox_sans_bold() {
    let (w1,h1) = get_text_bbox("hallo", 18,FontGroup::SansBold);
    let (w2,h2) = get_text_bbox("hallo", 24,FontGroup::SansBold);
    let (w3,h3) = get_text_bbox("hallo", 36,FontGroup::SansBold);
    assert_eq!((w1, h1),(41, 18));
    assert_eq!((w2, h2),(54, 24));
    assert_eq!((w3, h3),(81, 36));
}

#[test]
fn test_get_text_bbox_sans_bold_italic() {
    let (w1,h1) = get_text_bbox("hallo", 18,FontGroup::SansBoldItalic);
    let (w2,h2) = get_text_bbox("hallo", 24,FontGroup::SansBoldItalic);
    let (w3,h3) = get_text_bbox("hallo", 36,FontGroup::SansBoldItalic);
    assert_eq!((w1, h1),(42, 18));
    assert_eq!((w2, h2),(56, 24));
    assert_eq!((w3, h3),(83, 36));
}

#[test]
fn test_split_text_in_lines() {
    let text = "hallo";
    let lines = split_text_in_lines(text, 40, 18,FontGroup::Sans);
    assert_eq!(lines.len(), 1);

    let lines = split_text_in_lines(text, 20, 18,FontGroup::Sans);
    assert_eq!(lines.len(), 1);

    let text2 = "hallo hallo";
    let lines = split_text_in_lines(text2, 20, 18,FontGroup::Sans);
    assert_eq!(lines.len(), 2);

    let text3 = "hallo hallo hallo";
    let lines = split_text_in_lines(text3, 20, 18,FontGroup::Sans);
    assert_eq!(lines.len(), 3);
}

