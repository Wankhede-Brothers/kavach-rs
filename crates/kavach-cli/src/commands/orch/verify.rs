//! Orch verify: post-implementation verification checklist.
//! Checks that all quality gates pass before marking task complete.

use std::io::Write;

use crate::session;

pub fn run() -> anyhow::Result<()> {
    let sess = session::get_or_create_session();
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[ORCH:VERIFY]")?;
    writeln!(w, "project: {}", sess.project)?;
    writeln!(w, "date: {}", sess.today)?;
    writeln!(w)?;

    let mut pass = true;

    writeln!(w, "[CHECKLIST]")?;

    // Research done?
    let research_ok = sess.research_done;
    writeln!(
        w,
        "  research: {}",
        if research_ok { "PASS" } else { "FAIL" }
    )?;
    if !research_ok {
        pass = false;
    }

    // Aegis verified?
    let aegis_ok = sess.aegis_verified;
    writeln!(w, "  aegis: {}", if aegis_ok { "PASS" } else { "PENDING" })?;
    if !aegis_ok {
        pass = false;
    }

    // Task status
    let task_ok = sess.task_status == "completed" || sess.task_status == "landing";
    writeln!(
        w,
        "  task: {}",
        if sess.has_task() {
            if task_ok {
                "PASS"
            } else {
                "IN_PROGRESS"
            }
        } else {
            "N/A"
        }
    )?;

    // Files modified
    writeln!(w, "  files_modified: {}", sess.files_modified.len())?;

    writeln!(w)?;
    if pass {
        writeln!(w, "[RESULT] VERIFIED")?;
    } else {
        writeln!(w, "[RESULT] INCOMPLETE")?;
        if !research_ok {
            writeln!(w, "  action: Run WebSearch before implementation")?;
        }
        if !aegis_ok {
            writeln!(w, "  action: Run kavach orch aegis")?;
        }
    }

    Ok(())
}
