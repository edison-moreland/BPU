module toggle_reset(
    input  logic t, r,
    output logic q,
);
    always_ff @(posedge t, posedge r) begin
        if (r)
            q <= 0;
        else if(t)
            q <= ~q;
    end
endmodule