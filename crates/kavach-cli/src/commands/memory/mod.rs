pub mod bank;
pub mod write;
pub mod rpc;
pub mod stm;
pub mod kanban;
pub mod sync;
pub mod view;

use clap::Subcommand;
use super::cli_print_fmt;

#[derive(Subcommand)]
pub enum MemoryCommand {
    Bank(bank::BankArgs),
    Write(write::WriteArgs),
    /// JSON-RPC on stdin
    Rpc,
    Stm,
    Inject,
    Spec,
    Kanban(kanban::KanbanArgs),
    View,
    Sync(sync::SyncArgs),
}

pub fn dispatch(cmd: MemoryCommand) -> anyhow::Result<()> {
    match cmd {
        MemoryCommand::Bank(args) => bank::run(args),
        MemoryCommand::Write(args) => write::run(args),
        MemoryCommand::Rpc => rpc::run(),
        MemoryCommand::Stm => { cli_print_fmt("[STUB] memory stm".into()); Ok(()) }
        MemoryCommand::Inject => { cli_print_fmt("[STUB] memory inject".into()); Ok(()) }
        MemoryCommand::Spec => { cli_print_fmt("[STUB] memory spec".into()); Ok(()) }
        MemoryCommand::Kanban(args) => kanban::run(args),
        MemoryCommand::View => view::run(),
        MemoryCommand::Sync(args) => sync::run(args),
    }
}
