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

fn get_font_file(font_group: FontGroup) -> Vec<u8> {
    let cur_dir = env::current_dir().unwrap();
    let cur_dir_str = cur_dir.to_str().unwrap();

    let font_file = match font_group {
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

/// Returns (width, height) of a given text
/// ```rust
/// use bimpl::text_helper::{get_text_bbox, FontGroup};
/// 
/// let (w1,h1) = get_text_bbox("hallo", 18,FontGroup::Mono);
/// assert_eq!((w1, h1),(53, 18));
/// ```
pub fn get_text_bbox(text: &str, font_size: usize, font_group: FontGroup) -> (u32, u32) {
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
