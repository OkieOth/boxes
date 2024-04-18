use crate::layouts;
use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Debug)]
pub struct Doc {
    pub layout: layouts::Layout,
}
