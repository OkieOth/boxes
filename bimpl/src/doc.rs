use crate::layouts;
use serde::Deserialize;

#[derive(Deserialize, Debug)]
pub struct Doc {
    pub layout: layouts::Layout,
}
