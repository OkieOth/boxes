use crate::layouts;
use serde::Deserialize;

#[derive(Deserialize)]
pub struct Doc {
    pub layout: layouts::Layout,
}
