use crate::layouts;
use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Debug)]
pub struct DocConfig {
    pub max_default_with: u32,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Doc {
    pub layout: layouts::Layout,
    pub config: Option<DocConfig>
}
